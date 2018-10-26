package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	uuid "github.com/google/uuid"
	"github.com/pkg/errors"
)

const maxPlayers = 2

var (
	errRoomIsFull = errors.New("Room is full")
)

type RoomParameters struct {
	GameDuration time.Duration
}

type Room struct {
	game      *Game
	gameLogic GameLogic

	ID string

	sync.RWMutex
	availableSlots int
	players        []*Player
	parameters     RoomParameters
	register       chan *Player
	isFullChan     chan struct{}
	isDoneChan     chan struct{}
	broadcastsIn   [maxPlayers]chan *Message
	broadcastsOut  [maxPlayers]chan *Message
}

func createRoomID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}

func NewRoom(g *Game, params RoomParameters) *Room {
	r := &Room{
		game:       g,
		ID:         createRoomID(),
		players:    make([]*Player, 0, maxPlayers),
		parameters: params,
		register:   make(chan *Player),
		isFullChan: make(chan struct{}),
		broadcastsIn: [maxPlayers]chan *Message{
			make(chan *Message),
			make(chan *Message),
		},
		broadcastsOut: [maxPlayers]chan *Message{
			make(chan *Message),
			make(chan *Message),
		},
		availableSlots: maxPlayers,
	}

	go r.listen()
	return r
}

func (r *Room) TakeSlot() {
	r.Lock()
	r.availableSlots--
	r.Unlock()
}

func (r *Room) IsFull() bool {
	r.RLock()
	defer r.RUnlock()
	return r.availableSlots == 0
}

func (r *Room) isAllPlayersCreated() bool {
	r.RLock()
	defer r.RUnlock()
	return len(r.players) == maxPlayers
}

func (r *Room) AddPlayer(player *Player) {
	r.register <- player
}

func (r *Room) startGame() {

	singletoneLogger.LogMessage(fmt.Sprintf("RoomID: %v, is starting game with:\nPlayer1: %v\nPlayer2",
		r.ID, r.players[0].playerData, r.players[1].playerData))

	for idx, p := range r.players {
		p.Start(r.broadcastsIn[idx], r.broadcastsOut[idx])
	}

	for {
		select {
		case firstMsg := <-r.broadcastsIn[0]:
			if firstMsg.msgType == TurnMsg {
				turn := &Turn{}
				err := firstMsg.ToUnmarshalData(turn)
				if err != nil {
					log.Println(err)
					break
				}
				r.gameLogic.FirstPlayerTurn(*turn)
			}
		case secondMsg := <-r.broadcastsIn[1]:
			if secondMsg.msgType == TurnMsg {
				turn := &Turn{}
				err := secondMsg.ToUnmarshalData(turn)
				if err != nil {
					log.Println(err)
					break
				}
				r.gameLogic.SecondPlayerTurn(*turn)
			}
		}
	}
}

func (r *Room) listen() {
LOOP:
	for {
		select {
		case p := <-r.register:
			if r.isAllPlayersCreated() {
				r.isFullChan <- struct{}{}
			} else {
				r.players = append(r.players, p)
				singletoneLogger.LogMessage(fmt.Sprintf("RoomID %v: Player was added: %v", p.playerData.Name))
			}
		case <-r.isFullChan:
			go r.startGame()
			break LOOP
		}
	}

	<-r.isDoneChan
	r.game.CloseRoom(r.ID)
}
