package chat

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gamelogic"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/randomGenerator"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
)

const maxPlayers = 100

var (
	errRoomIsFull          = errors.New("Room is full")
	errInvalidGameDuration = errors.New("Invalid game duration")
	errCloseWhileSearch    = errors.New("Room was closed in search")
	errInternalServerError = errors.New("Internal server error")
)

const (
	closeNormalGameOverFrame = "Game over"
)

var (
	startMsgWhite      *Message
	startMsgBlack      *Message
	infoMsgGameStarted *Message
	infoMsgAddedToRoom *Message

	errMsgInvalidJSON      *Message
	errMsgUnknownMsgType   *Message
	infoMsgRivalDisconnect *Message
)

func init() {
	infoMsgGameStarted, _ = MarshalToMessage(InfoMsg, &InfoMessage{"Game is started"})
	infoMsgAddedToRoom, _ = MarshalToMessage(InfoMsg, &InfoMessage{"Added to room"})
	errMsgInvalidJSON, _ = MarshalToMessage(ErrorMsg, &ErrorMessage{"Invalid JSON"})
	errMsgUnknownMsgType, _ = MarshalToMessage(ErrorMsg, &ErrorMessage{"Unknown message type"})
	// // errMsgInternalError, _ = MarshalToMessage(ErrorMsg, &ErrorMessage{"Write Internal server error"})
}

//easyjson:json
type RoomParameters struct {
	Duration int `json:"duration" example:"60"`
}

type Room struct {
	gameLogic gamelogic.GameLogic

	ID        string
	closeTime int // Через какое количество секунд, комната закроется если не заполнилась

	sync.RWMutex
	availableSlots int
	slotsGuids     []string

	players        []*User
	parameters     RoomParameters
	register       chan *User
	isFullChan     chan struct{}
	isDoneChan     chan struct{}
	broadcastAllIn chan *UserMsg
	broadcastsIn   []chan *UserMsg
	broadcastsOut  []chan *UserMsg
	closeErrors    [maxPlayers]chan error

	isFirstWinner bool
}

func NewRoom() *Room {
	id, err := randomGenerator.RandomUUID()
	if err != nil {
		singletoneLogger.LogError(err)
	}

	r := &Room{
		ID:            id,
		players:       make([]*User, 0, maxPlayers),
		register:      make(chan *User),
		isFullChan:    make(chan struct{}, 1),
		isDoneChan:    make(chan struct{}, 1),
		broadcastsOut: make([]chan *UserMsg, 0, maxPlayers),
		broadcastsIn:  make([]chan *UserMsg, 0, maxPlayers),
		closeErrors: [maxPlayers]chan error{
			make(chan error),
			make(chan error),
		},
		broadcastAllIn: make(chan *UserMsg, 100),
		availableSlots: maxPlayers,
		slotsGuids:     make([]string, 0, maxPlayers),
	}

	go r.listen()
	return r
}

func (r *Room) listen() {
	closeTimer := time.NewTimer(time.Second * time.Duration(r.closeTime))
	startTimer := time.Now()

	go r.listenUsers()
	for {
		select {
		case p := <-r.register:

			if !closeTimer.Stop() {
				<-closeTimer.C
			}
			closeTimer.Reset(time.Second * time.Duration(r.closeTime))
			singletoneLogger.LogMessage("Timer was reseted")

			singletoneLogger.LogMessage(fmt.Sprintf("RoomID %v: Player was added: %v", r.ID, p.UserData.Name))

			r.broadcastsIn = append(r.broadcastsIn, make(chan *UserMsg))
			r.broadcastsOut = append(r.broadcastsOut, make(chan *UserMsg))

			p.Start(r.broadcastsIn[len(r.broadcastsIn)-1], r.broadcastsOut[len(r.broadcastsOut)-1], nil)

			go func() {
				for {
					r.broadcastAllIn <- <-r.broadcastsIn[len(r.broadcastsIn)-1]
				}
			}()
		}
	}
	finishTime := time.Now().Sub(startTimer)
	singletoneLogger.LogMessage(fmt.Sprintf("%v", finishTime.Seconds()))
}

func (r *Room) resolveMsg(msg *UserMsg) {
	um := &UserMessage{}
	err := msg.Msg.UnmarshalData(um)
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}

	if um.To == "" {
		r.messageToAll(msg)
	}
}

func (r *Room) listenUsers() {
	for {
		select {
		case msg := <-r.broadcastAllIn:
			singletoneLogger.LogMessage(msg.Guid)
			switch msg.Msg.MsgType {
			case "msg":
				r.resolveMsg(msg)
			default:
				singletoneLogger.LogMessage("Unknown msg")
			}
		}
	}
}

func (r *Room) TakeSlot(guid string) {
	r.Lock()
	r.availableSlots--
	r.slotsGuids = append(r.slotsGuids, guid)
	r.Unlock()
}

func (r *Room) IsUserInRoomByGuid(guid string) bool {
	r.RLock()
	defer r.RUnlock()
	for _, g := range r.slotsGuids {
		if g == guid {
			return true
		}
	}

	return false
}

func (r *Room) IsFull() bool {
	r.RLock()
	defer r.RUnlock()
	return r.availableSlots == 0
}

func (r *Room) AddPlayer(player *User) {
	r.register <- player
}

func (r *Room) messageToAll(msg *UserMsg) {
	for _, out := range r.broadcastsOut {
		out <- msg
	}
}

func (r *Room) closeConnections(code int, msg string) {
	for _, p := range r.players {
		p.Close(code, msg)
	}
}
