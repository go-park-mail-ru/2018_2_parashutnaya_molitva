package game

import (
	"fmt"
	GRPCCore "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC/core"
	"golang.org/x/net/context"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type UserController interface {
	AddScoreToUser(ctx context.Context, scoreAdd *GRPCCore.ScoreAdd) (*GRPCCore.Nothing, error)
}

type PlayerData struct {
	Name           string
	Guid           string
	Score          int
}

type Player struct {
	conn       *websocket.Conn
	playerData PlayerData

	sync.RWMutex
	isClosed bool

	closeChan    chan struct{}
	isClosedChan bool
}

func NewPlayer(name, guid string, score int, conn *websocket.Conn) *Player {
	return &Player{
		conn: conn,
		playerData: PlayerData{
			name,
			guid,
			score,
		},
		isClosed:     false,
		isClosedChan: false,
		closeChan:    make(chan struct{}, 1),
	}
}

func (p *Player) GetGUID() string {
	return p.playerData.Guid
}

func (p *Player) GetName() string {
	return p.playerData.Name
}

func (p *Player) GetScore() int {
	return p.playerData.Score
}

func (p *Player) IsClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.isClosed
}

func (p *Player) Send(msg *Message) error {
	err := p.conn.WriteJSON(msg)
	if websocket.IsUnexpectedCloseError(errors.WithStack(err)) {
		p.close()
	}

	return errors.WithStack(err)
}

func (p *Player) close() {
	p.Lock()
	p.isClosed = true
	p.Unlock()
}

func (p *Player) closeByChan() {
	if p.isClosedChan {
		return
	}
	p.isClosedChan = true
	close(p.closeChan)
}

// Close - закрывает соедениение с номером ошибки и причиной
// Безопасен для уже закрытых соеднинений
func (p *Player) Close(code int, msg string) {
	if p.IsClosed() {
		return
	}
	sendCloseError(p.conn, code, msg)
	p.close()
	p.closeByChan()
	singletoneLogger.LogMessage("Close by chan")
	p.conn.Close()

}

// Читает сообщения из канал и записывает в соединение
func (p *Player) write(out <-chan *Message, chanCloseError chan<- error) {
	pingTicker := time.NewTicker(time.Second * time.Duration(gameConfig.PingPeriod))
	defer pingTicker.Stop()
	for {
		select {
		case msg, ok := <-out:
			if !ok {
				return
			}

			sendMsg, err := msg.MarshalJSON()
			if err != nil {
				singletoneLogger.LogError(errInavlidMsgFormat)
				continue
			}

			err = p.conn.WriteMessage(websocket.TextMessage, sendMsg)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseTryAgainLater) {
				p.close()
				p.closeByChan()
				return
			} else if err != nil {
				// Не блокируемся на ошибики, если уже закрыто соединение
				p.close()
				select {
				case chanCloseError <- errors.WithStack(err):
				case <-p.closeChan:
				}
				// В следующий раз как будет ошибка, но чтение еще не закончилось, chanCloseError не заблочиться
				p.closeByChan()
				return
			}
		case <-p.closeChan:
			return
		case <-pingTicker.C:
			err := p.conn.WriteMessage(websocket.PingMessage, nil)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseTryAgainLater) {
				p.close()
				p.closeByChan()
				return
			} else if err != nil {
				// Не блокируемся на ошибики, если уже закрыто соединение
				p.close()
				select {
				case chanCloseError <- errors.WithStack(err):
				case <-p.closeChan:
				}
				// В следующий раз как будет ошибка, но чтение еще не закончилось, chanCloseError не заблочиться
				p.closeByChan()
				return
			}
		}
	}
}

// Читает сообщения из соединения и записывает в канал
func (p *Player) read(in chan<- *Message, chanCloseError chan<- error) {
	readMessage := &Message{}

	p.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(gameConfig.PongWait)))
	p.conn.SetPongHandler(func(string) error {
		p.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(gameConfig.PongWait)))
		return nil
	})

	for !p.IsClosed() {

		_, msg, err := p.conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseTryAgainLater) {
			p.close()
			p.closeByChan()
			return
		} else if err != nil {
			p.close()
			// Не блокируемся на ошибики, если уже закрыто соединение
			select {
			case chanCloseError <- errors.WithStack(err):
			case <-p.closeChan:
			}
			// В следующий раз как будет ошибка, но чтение еще не закончилось, chanCloseError не заблочиться
			p.closeByChan()
			return
		}

		readMessage, err = UnmarshalToMessage(msg)
		if err != nil {
			singletoneLogger.LogError(errors.WithStack(err))
			sendError(p.conn, errInavlidMsgFormat.Error())
			continue
		}

		in <- readMessage
	}
}

func (p *Player) Start(in chan<- *Message, out <-chan *Message, chanCloseError chan<- error) {
	singletoneLogger.LogMessage(fmt.Sprintf("Start listen player: %v", p.playerData.Name))

	go p.read(in, chanCloseError)
	go p.write(out, chanCloseError)
}
