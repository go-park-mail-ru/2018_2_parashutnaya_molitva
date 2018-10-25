package game

import (
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

var (
	errNoRoom = errors.New("No Room")
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

	go g.Listen()
	return g
}

func (g *Game) InitRoom(roomParameters RoomParameters) string {
	g.mx.RLock()
	for id, room := range g.rooms {
		if room.parameters == roomParameters && !room.IsFull() {
			return id
		}
	}
	g.mx.RUnlock()

	r := NewRoom(g, roomParameters)
	g.createRoom <- r
	return r.ID
}

func (g *Game) InitConnection(roomID, name string, score int, conn *websocket.Conn) error {
	g.mx.RLock()
	room, exist := g.rooms[roomID]
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

func (g *Game) Listen() {
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
		singletoneLogger.LogMessage(fmt.Sprintf("RoomId: %v", id))
		singletoneLogger.LogMessage(fmt.Sprintf("RoomParametrs: %#v", r.parameters.GameDuration.String()))
		singletoneLogger.LogMessage("------------")
	}
	g.mx.RUnlock()
}
