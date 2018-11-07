package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type UserStorage interface {
	AddScore(int) error
}

type PlayerData struct {
	Name  string
	Score int
	User  UserStorage
}

type Player struct {
	conn       *websocket.Conn
	playerData PlayerData

	sync.RWMutex
	isClosed bool

	closeChan chan struct{}
}

func NewPlayer(name string, score int, user UserStorage, conn *websocket.Conn) *Player {
	return &Player{
		conn: conn,
		playerData: PlayerData{
			name,
			score,
			user,
		},
		isClosed:  false,
		closeChan: make(chan struct{}, 1),
	}
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
	if websocket.IsUnexpectedCloseError(err) {
		p.close()
	}

	return errors.WithStack(err)
}

func (p *Player) close() {
	p.Lock()
	p.isClosed = true
	p.Unlock()
	p.closeChan <- struct{}{}
}

// Close - закрывает соедениение с номером ошибки и причиной
// Безопасен для уже закрытых соеднинений
func (p *Player) Close(code int, msg string) {
	if p.IsClosed() {
		return
	}
	p.close()
	sendCloseError(p.conn, code, msg)
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
			if err != nil {
				p.close()
				chanCloseError <- errors.WithStack(err)
				return
			}
		case <-p.closeChan:
			return
		case <-pingTicker.C:
			err := p.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				p.close()
				chanCloseError <- errors.WithStack(err)
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
		if err != nil {
			p.close()
			chanCloseError <- errors.WithStack(err)
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
