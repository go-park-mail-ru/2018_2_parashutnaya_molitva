package game

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC/core"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"

	"github.com/gorilla/websocket"
)

//go:generate easyjson -pkg

type Game struct {
	GRPCCore core.CoreClient

	mx             sync.RWMutex
	rooms          map[string]*Room
	searchingGuids map[string]struct{}

	closeRoom  chan string
	createRoom chan *Room
}

func NewGame(core core.CoreClient) *Game {
	g := &Game{
		GRPCCore:       core,
		mx:             sync.RWMutex{},
		rooms:          map[string]*Room{},
		closeRoom:      make(chan string),
		createRoom:     make(chan *Room),
		searchingGuids: make(map[string]struct{}),
	}

	go g.listen()
	return g
}

func (g *Game) findRoom(guid string, roomParameters RoomParameters) (string, error) {
	g.mx.RLock()
	for id, room := range g.rooms {
		g.mx.RUnlock()
		if room.IsUserInRoomByGuid(guid) {
			return "", errAlreadyInRoom
		}

		if room.parameters == roomParameters && !room.IsFull() {
			room.TakeSlot(guid)
			singletoneLogger.LogMessage(fmt.Sprintf("searching: %#v, ", g.searchingGuids))
			g.deleteFromSearching(guid)

			return id, nil
		}
		g.mx.RLock()
	}
	g.mx.RUnlock()
	return "", nil
}

func (g *Game) isUserAlreadySearching(guid string) bool {
	g.mx.RLock()
	defer g.mx.RUnlock()
	_, isExist := g.searchingGuids[guid]
	return isExist
}

func (g *Game) addUserInSearching(guid string) {
	g.mx.Lock()
	g.searchingGuids[guid] = struct{}{}
	g.mx.Unlock()
}

func (g *Game) deleteFromSearching(guid string) {
	g.mx.Lock()
	delete(g.searchingGuids, guid)
	g.mx.Unlock()
}

func (g *Game) InitRoom(guid string, roomParameters RoomParameters) (string, error) {

	singletoneLogger.LogMessage(fmt.Sprintf("In searching: %#v", g.searchingGuids))
	if g.isUserAlreadySearching(guid) {
		return "", errAlreadySearching
	}
	g.addUserInSearching(guid)
	defer g.deleteFromSearching(guid)
	id, err := g.findRoom(guid, roomParameters)
	singletoneLogger.LogMessage(fmt.Sprintf("In searching: %#v", g.searchingGuids))
	if id != "" && err == nil {
		return id, nil
	} else if err != nil {
		return "", err
	}

	r := NewRoom(g, roomParameters, gameConfig.CloseRoomDeadline)
	r.TakeSlot(guid)
	g.createRoom <- r
	return r.ID, nil
}

func (g *Game) readInitMessage(done chan struct{}, conn *websocket.Conn) (chan *InitMessage, chan error) {
	chanMessage := make(chan *InitMessage, 1)
	chanCloseError := make(chan error)
	go func() {
		for {

			select {
			case <-done:
				return
			default:
			}

			_, rawMsg, err := conn.ReadMessage()
			if websocket.IsUnexpectedCloseError(err) {
				singletoneLogger.LogError(err)
				chanCloseError <- err
				return
			} else if err != nil {
				sendError(conn, err.Error())
				continue
			}

			if rawMsg == nil {
				continue
			}

			singletoneLogger.LogMessage(string(rawMsg))

			msg, err := UnmarshalToMessage(rawMsg)
			if err != nil {
				singletoneLogger.LogError(err)
				sendError(conn, errInavlidMsgFormat.Error())
				continue
			}

			if msg.MsgType != InitMsg {
				sendError(conn, errInvalidMsgTypeInit.Error())
				continue
			}

			initMsg := &InitMessage{}
			err = msg.UnmarshalData(initMsg)
			if err != nil {
				sendError(conn, errInvalidMsgInitFormat.Error())
				continue
			}

			chanMessage <- initMsg
			return
		}
	}()

	return chanMessage, chanCloseError
}

func (g *Game) initConnection(name, guid string, score int, conn *websocket.Conn) {
	done := make(chan struct{})
	t := time.NewTimer(time.Second * time.Duration(gameConfig.InitMessageDeadline))
	initChan, closeErrorChan := g.readInitMessage(done, conn)
	select {
	case initMessage := <-initChan:
		close(done)
		g.mx.RLock()
		room, exist := g.rooms[initMessage.RoomId]
		g.mx.RUnlock()

		if !exist {
			singletoneLogger.LogError(errNoRoom)
			return
		}

		room.AddPlayer(NewPlayer(name, guid, score, conn))
		singletoneLogger.LogMessage(fmt.Sprintf("Successfully init msg was read: %v", initMessage))
	case <-t.C:
		close(done)

		singletoneLogger.LogError(errInitMsgWaitTooLong)
		sendCloseError(conn, websocket.CloseNormalClosure, errInitMsgWaitTooLong.Error())
	case closeError := <-closeErrorChan:
		close(done)

		if !t.Stop() {
			<-t.C
		}

		singletoneLogger.LogError(closeError)
	}
}

func (g *Game) InitConnection(name, guid string, score int, conn *websocket.Conn) {
	go g.initConnection(name, guid, score, conn)
}

func (g *Game) CloseRoom(roomID string) {
	g.closeRoom <- roomID
}

func (g *Game) listen() {
	singletoneLogger.LogMessage("Game server start listen")
	for {
		select {
		case roomID := <-g.closeRoom:

			g.mx.Lock()
			delete(g.rooms, roomID)
			g.mx.Unlock()
			singletoneLogger.LogMessage(fmt.Sprintf("Room: %v, was deleted", roomID))
			g.printGameState()
		case room := <-g.createRoom:

			g.mx.Lock()
			g.rooms[room.ID] = room
			g.mx.Unlock()

			g.printGameState()
		}
	}
}

func (g *Game) printGameState() {
	singletoneLogger.LogMessage("------------")
	singletoneLogger.LogMessage("Rooms")
	g.mx.RLock()
	for id, r := range g.rooms {
		g.mx.RUnlock()
		singletoneLogger.LogMessage(fmt.Sprintf("RoomId: %v", id))
		singletoneLogger.LogMessage(fmt.Sprintf("RoomParametrs: %#v", r.parameters.Duration))
		g.mx.RLock()
	}
	g.mx.RUnlock()
	singletoneLogger.LogMessage("------------")
}
