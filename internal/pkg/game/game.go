package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

var (
	errNoRoom               = errors.New("No Room")
	errUnexpectedClose      = errors.New("Unexpected Close")
	errInitMsgWaitTooLong   = errors.New("Waiting Init Message too long")
	errInvalidMsgTypeInit   = errors.New("Invalid type: wating for Init")
	errInvalidMsgInitFormat = errors.New("Invalid init message format")
	errInavlidMsgFormat     = errors.New("Invalid message format")
)

type Game struct {
	mx         sync.RWMutex
	rooms      map[string]*Room
	closeRoom  chan string
	createRoom chan *Room
}

func NewGame() *Game {
	g := &Game{
		mx:         sync.RWMutex{},
		rooms:      map[string]*Room{},
		closeRoom:  make(chan string),
		createRoom: make(chan *Room),
	}

	go g.listen()
	return g
}

func (g *Game) findRoom(roomParameters RoomParameters) string {
	g.mx.RLock()
	defer g.mx.RUnlock()
	for id, room := range g.rooms {
		if room.parameters == roomParameters && !room.IsFull() {
			room.TakeSlot()
			return id
		}
	}

	return ""
}

func (g *Game) InitRoom(roomParameters RoomParameters) string {
	id := g.findRoom(roomParameters)

	if id != "" {
		return id
	}

	r := NewRoom(g, roomParameters)
	r.TakeSlot()
	g.createRoom <- r
	return r.ID
}

func (g *Game) readInitMessage(done chan struct{}, conn *websocket.Conn) chan *InitMessage {
	chanMessage := make(chan *InitMessage, 1)
	go func() {
		for {

			select {
			case <-done:
				return
			default:
			}

			_, rawMsg, err := conn.ReadMessage()
			if websocket.IsUnexpectedCloseError(err) {
				conn.Close()
				singletoneLogger.LogError(err)
				return
			}

			if rawMsg == nil {
				continue
			}

			singletoneLogger.LogMessage(string(rawMsg))

			msg, err := UnmarshalToMessage(rawMsg)
			if err != nil {
				sendError(conn, errInavlidMsgFormat.Error())
				continue
			}
			singletoneLogger.LogMessage(fmt.Sprintf("%#v", msg))
			if msg.MsgType != InitMsg {
				sendError(conn, errInvalidMsgTypeInit.Error())
				// singletoneLogger.LogError(errInvalidMsgTypeInit)
				continue
			}

			initMsg := &InitMessage{}
			err = msg.ToUnmarshalData(initMsg)
			if err != nil {
				sendError(conn, errInvalidMsgInitFormat.Error())
				continue
			}

			chanMessage <- initMsg
			return
		}
	}()

	return chanMessage
}

func (g *Game) initConnection(name string, score int, conn *websocket.Conn) {
	done := make(chan struct{})
	t := time.NewTimer(time.Second * time.Duration(gameConfig.InitMessageDeadline)).C

	select {
	case initMessage := <-g.readInitMessage(done, conn):
		close(done)
		g.mx.RLock()
		room, exist := g.rooms[initMessage.RoomId]
		g.mx.RUnlock()

		if !exist {
			singletoneLogger.LogError(errNoRoom)
			return
		}

		room.AddPlayer(NewPlayer(name, score, conn))
		singletoneLogger.LogMessage(fmt.Sprintf("Successfully init msg was read: %v", initMessage))
	case <-t:
		close(done)
		sendError(conn, errInitMsgWaitTooLong.Error())
		singletoneLogger.LogError(errInitMsgWaitTooLong)
		conn.Close()
	}
}

func (g *Game) InitConnection(name string, score int, conn *websocket.Conn) {
	go g.initConnection(name, score, conn)
}

func (g *Game) CloseRoom(roomID string) {
	g.closeRoom <- roomID
}

func (g *Game) listen() {
	singletoneLogger.LogMessage("Game server start listen")
	for {
		select {
		case roomID := <-g.closeRoom:
			delete(g.rooms, roomID)
			g.printGameState()
		case room := <-g.createRoom:
			g.rooms[room.ID] = room
			g.printGameState()
		}
	}
}

func (g *Game) printGameState() {
	g.mx.RLock()
	for id, r := range g.rooms {
		g.mx.RUnlock()
		singletoneLogger.LogMessage(fmt.Sprintf("RoomId: %v", id))
		singletoneLogger.LogMessage(fmt.Sprintf("RoomParametrs: %#v", r.parameters.Duration))
		singletoneLogger.LogMessage("------------")
		g.mx.RLock()
	}
	g.mx.RUnlock()
}
