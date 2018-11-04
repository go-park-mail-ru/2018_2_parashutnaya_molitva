package game

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	uuid "github.com/google/uuid"
	"github.com/pkg/errors"
)

const maxPlayers = 2

var (
	errRoomIsFull          = errors.New("Room is full")
	errInvalidGameDuration = errors.New("Invalid game duration")
)

//easyjson:json
type RoomParameters struct {
	Duration int `json:"duration" example:"60"`
}

func (r *RoomParameters) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *RoomParameters) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(r)
	return data, err
}

func (r *RoomParameters) Validate() (string, error) {
	if err := ValidateDuration(r.Duration); err != nil {
		return "duration", err
	}

	return "", nil
}

func ValidateDuration(duration int) error {
	for _, d := range gameConfig.ValidGameDuration {
		if d == duration {
			return nil
		}
	}

	return errInvalidGameDuration
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

	gameIsStartedMsg, _ := MarshalToMessage(InfoMsg, &InfoMessage{"Game is started"})
	for _, out := range r.broadcastsOut {
		out <- gameIsStartedMsg
	}
	for {
		select {
		case firstMsg := <-r.broadcastsIn[0]:
			if firstMsg.MsgType == TurnMsg {
				turn := &Turn{}
				err := firstMsg.ToUnmarshalData(turn)
				if err != nil {
					log.Println(err)
					break
				}
				r.gameLogic.FirstPlayerTurn(*turn)
			}
		case secondMsg := <-r.broadcastsIn[1]:
			if secondMsg.MsgType == TurnMsg {
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
				
				msg, _ := MarshalToMessage(InfoMsg, &InfoMessage{"Added to room"})
				p.Send(msg)
			}
		case <-r.isFullChan:
			go r.startGame()
			break LOOP
		}
	}

	<-r.isDoneChan
	r.game.CloseRoom(r.ID)
}
