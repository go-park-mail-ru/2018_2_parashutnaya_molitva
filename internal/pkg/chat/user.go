package chat

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

type UserData struct {
	Name string
	Guid string
}

type UserMsg struct {
	Guid string
	Msg  *Message
}

type User struct {
	conn     *websocket.Conn
	UserData UserData

	sync.RWMutex
	isClosed bool

	closeChan    chan struct{}
	isClosedChan bool
}

func NewUser(name, guid string, conn *websocket.Conn) *User {
	return &User{
		conn: conn,
		UserData: UserData{
			name,
			guid,
		},
		isClosed:     false,
		isClosedChan: false,
		closeChan:    make(chan struct{}, 1),
	}
}

func (p *User) GetGUID() string {
	return p.UserData.Guid
}

func (p *User) GetName() string {
	return p.UserData.Name
}

func (p *User) IsClosed() bool {
	p.RLock()
	defer p.RUnlock()
	return p.isClosed
}

func (p *User) Send(msg *Message) error {
	err := p.conn.WriteJSON(msg)
	if websocket.IsUnexpectedCloseError(errors.WithStack(err)) {
		p.close()
	}

	return errors.WithStack(err)
}

func (p *User) close() {
	p.Lock()
	p.isClosed = true
	p.Unlock()
}

func (p *User) closeByChan() {
	if p.isClosedChan {
		return
	}
	p.isClosedChan = true
	close(p.closeChan)
}

// Close - закрывает соедениение с номером ошибки и причиной
// Безопасен для уже закрытых соеднинений
func (p *User) Close(code int, msg string) {
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
func (p *User) write(out <-chan *UserMsg, chanCloseError chan<- error) {
	pingTicker := time.NewTicker(time.Second * 10)
	defer pingTicker.Stop()
	for {
		select {
		case msg, ok := <-out:
			if !ok {
				return
			}
			if msg.Guid == p.UserData.Guid {
				singletoneLogger.LogMessage("Same")
				continue
			}

			sendMsg, err := msg.Msg.MarshalJSON()
			if err != nil {
				singletoneLogger.LogError(errInavlidMsgFormat)
				continue
			}

			err = p.conn.WriteMessage(websocket.TextMessage, sendMsg)
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseTryAgainLater) {
			} else if err != nil {
				return
			}
		}
	}
}

// Читает сообщения из соединения и записывает в канал
func (p *User) read(in chan<- *UserMsg, chanCloseError chan<- error) {
	readMessage := &Message{}

	p.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	p.conn.SetPongHandler(func(string) error {
		p.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
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

		singletoneLogger.LogMessage(string(msg))
		readMessage, err = UnmarshalToMessage(msg)
		if err != nil {
			singletoneLogger.LogError(errors.WithStack(err))
			sendError(p.conn, errInavlidMsgFormat.Error())
			continue
		}
		singletoneLogger.LogMessage(string(msg))
		in <- &UserMsg{Guid: p.UserData.Guid, Msg: readMessage}
	}
}

func (p *User) Start(in chan<- *UserMsg, out <-chan *UserMsg, chanCloseError chan<- error) {
	singletoneLogger.LogMessage(fmt.Sprintf("Start listen userController: %v", p.UserData.Name))

	go p.read(in, chanCloseError)
	go p.write(out, chanCloseError)
}
