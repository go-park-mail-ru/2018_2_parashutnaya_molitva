package game

import (
	"log"
	"time"

	uuid "github.com/google/uuid"
)

const maxPlayers = 2

type RoomParameters struct {
	GameDuration time.Duration
}

type Room struct {
	game      *Game
	gameLogic GameLogic

	ID            string
	players       []*Player
	parameters    RoomParameters
	register      chan *Player
	isFullChan    chan struct{}
	isDoneChan    chan struct{}
	broadcastsIn  [maxPlayers]chan *Message
	broadcastsOut [maxPlayers]chan *Message
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
	}

	go r.listen()
	return r
}

func (r *Room) IsFull() bool {
	return len(r.players) == maxPlayers
}

func (r *Room) AddPlayer(player *Player) {
	r.register <- player
}

func (r *Room) startGame() {
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
			if len(r.players) >= maxPlayers {
				r.isFullChan <- struct{}{}
			} else {
				r.players = append(r.players, p)
			}
		case <-r.isFullChan:
			go r.startGame()
			break LOOP
		}
	}

	<-r.isDoneChan
	r.game.CloseRoom(r.ID)
}
