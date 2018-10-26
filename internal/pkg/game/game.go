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
	errNoRoom             = errors.New("No Room")
	errUnexpectedClose    = errors.New("Unexpected Close")
	errInitMsgWaitTooLong = errors.New("Waiting Init Message too long")
)

type Game struct {
	mx         sync.RWMutex
	rooms      map[string]*Room
	closeRoom  chan string
	createRoom chan *Room
}

type InitMessage struct {
	RoomId string `json:"roomid"`
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

func (g *Game) readInitMessage(conn *websocket.Conn) (*InitMessage, error) {
	chanMessage := make(chan *InitMessage)
	closeErrMessage := make(chan error)
	t := time.NewTimer(time.Second * time.Duration(gameConfig.initMessageDeadline)).C
	go func() {
		for {
			_, rawMsg, err := conn.ReadMessage()
			if websocket.IsUnexpectedCloseError(err) {
				closeErrMessage <- err
			}
			msg, err := UnmarshalToMessage(rawMsg)
			if err != nil {
				singletoneLogger.LogError(err)
				continue
			}

			initMsg := &InitMessage{}
			err = msg.ToUnmarshalData(initMsg)
			if err != nil {
				singletoneLogger.LogError(err)
				continue
			}
			chanMessage <- initMsg
		}
	}()

	select {
	case initMessage := <-chanMessage:
		return initMessage, nil
	case err := <-closeErrMessage:
		return nil, err
	case <-t:
		return nil, errInitMsgWaitTooLong
	}
}

func (g *Game) InitConnection(name string, score int, conn *websocket.Conn) error {
	initMessage, err := g.readInitMessage(conn)
	if err != nil {
		return err
	}

	singletoneLogger.LogMessage(fmt.Sprintf("Successfully init msg was read: %v", initMessage))

	g.mx.RLock()
	room, exist := g.rooms[initMessage.RoomId]
	g.mx.RUnlock()
	if !exist {
		return errNoRoom
	}

	room.AddPlayer(NewPlayer(name, score, conn))
	return nil
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
		singletoneLogger.LogMessage(fmt.Sprintf("RoomParametrs: %#v", r.parameters.GameDuration.String()))
		singletoneLogger.LogMessage("------------")
		g.mx.RLock()
	}
	g.mx.RUnlock()
}
