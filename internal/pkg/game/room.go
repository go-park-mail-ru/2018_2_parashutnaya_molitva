package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gamelogic"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/randomGenerator"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const maxPlayers = 2
const scoreFactor = 5

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
	startMsgWhite, _ = MarshalToMessage(StartMsg, &StartMessage{"w"})
	startMsgBlack, _ = MarshalToMessage(StartMsg, &StartMessage{"b"})
	infoMsgGameStarted, _ = MarshalToMessage(InfoMsg, &InfoMessage{"Game is started"})
	infoMsgAddedToRoom, _ = MarshalToMessage(InfoMsg, &InfoMessage{"Added to room"})
	errMsgInvalidJSON, _ = MarshalToMessage(ErrorMsg, &ErrorMessage{"Invalid JSON"})
	errMsgUnknownMsgType, _ = MarshalToMessage(ErrorMsg, &ErrorMessage{"Unknown message type"})
	// errMsgInternalError, _ = MarshalToMessage(ErrorMsg, &ErrorMessage{"Write Internal server error"})
	infoMsgRivalDisconnect, _ = MarshalToMessage(InfoMsg, &InfoMessage{"Rival was disconnected"})

}

//easyjson:json
type RoomParameters struct {
	Duration int `json:"duration" example:"60"`
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
	gameLogic gamelogic.GameLogic

	ID        string
	closeTime int // Через какое количество секунд, комната закроется если не заполнилась

	sync.RWMutex
	availableSlots int
	slotsGuids     []string

	players       []*Player
	parameters    RoomParameters
	register      chan *Player
	isFullChan    chan struct{}
	isDoneChan    chan struct{}
	broadcastsIn  [maxPlayers]chan *Message
	broadcastsOut [maxPlayers]chan *Message
	closeErrors   [maxPlayers]chan error

	isFirstWinner bool
}

func NewRoom(g *Game, params RoomParameters, closeTime int) *Room {
	id, err := randomGenerator.RandomUUID()
	if err != nil {
		singletoneLogger.LogError(err)
	}

	r := &Room{
		game:       g,
		ID:         id,
		players:    make([]*Player, 0, maxPlayers),
		parameters: params,
		register:   make(chan *Player),
		isFullChan: make(chan struct{}, 1),
		isDoneChan: make(chan struct{}, 1),
		broadcastsIn: [maxPlayers]chan *Message{
			make(chan *Message),
			make(chan *Message),
		},
		broadcastsOut: [maxPlayers]chan *Message{
			make(chan *Message),
			make(chan *Message),
		},
		closeErrors: [maxPlayers]chan error{
			make(chan error),
			make(chan error),
		},
		availableSlots: maxPlayers,
		closeTime:      closeTime,
		gameLogic:      gamelogic.NewChessGameLogic(params.Duration),
		slotsGuids:     make([]string, 0, maxPlayers),
	}

	go r.listen()
	return r
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

func (r *Room) isAllPlayersCreated() bool {
	r.RLock()
	defer r.RUnlock()
	return len(r.players) == maxPlayers
}

func (r *Room) AddPlayer(player *Player) {
	r.register <- player
}

func (r *Room) messageToAll(msg *Message) {
	for _, out := range r.broadcastsOut {
		out <- msg
	}
}

func (r *Room) startGame() {

	singletoneLogger.LogMessage(fmt.Sprintf("RoomID: %v, is started game with:\nPlayer1: %v\nPlayer2: %v",
		r.ID, r.players[0].playerData, r.players[1].playerData))

	for idx, p := range r.players {
		p.Start(r.broadcastsIn[idx], r.broadcastsOut[idx], r.closeErrors[idx])
	}

	r.messageToAll(infoMsgGameStarted)

	isFirstWhite, gameOverChannel := r.gameLogic.Start()
	if isFirstWhite {
		singletoneLogger.LogMessage(fmt.Sprintf("Player1 is white: %v\nPlayer2 is black: %v",
			r.players[0].GetName(), r.players[1].GetName()))
		r.broadcastsOut[0] <- startMsgWhite
		r.broadcastsOut[1] <- startMsgBlack
	} else {
		singletoneLogger.LogMessage(fmt.Sprintf("Player1 is black: %v\nPlayer2 is white: %v",
			r.players[0].GetName(), r.players[1].GetName()))
		r.broadcastsOut[0] <- startMsgBlack
		r.broadcastsOut[1] <- startMsgWhite
	}

	gameIsEnded := false
	for !gameIsEnded {
		select {
		case firstMsg := <-r.broadcastsIn[0]:
			switch firstMsg.MsgType {
			case TurnMsg:
				r.turn(firstMsg, r.broadcastsOut[0], isFirstWhite)
			default:
				r.broadcastsOut[0] <- errMsgUnknownMsgType
			}
		case secondMsg := <-r.broadcastsIn[1]:
			switch secondMsg.MsgType {
			case TurnMsg:
				r.turn(secondMsg, r.broadcastsOut[1], !isFirstWhite)
			default:
				r.broadcastsOut[1] <- errMsgUnknownMsgType
			}
		case closeFirstErr := <-r.closeErrors[0]:
			singletoneLogger.LogError(closeFirstErr)
			r.broadcastsOut[1] <- infoMsgRivalDisconnect
			gameIsEnded = true

			r.endGame(r.players[1], r.players[0])
		case closeSecondErr := <-r.closeErrors[1]:
			singletoneLogger.LogError(closeSecondErr)
			r.broadcastsOut[0] <- infoMsgRivalDisconnect
			gameIsEnded = true

			r.endGame(r.players[0], r.players[1])
		case result := <-gameOverChannel:
			gameIsEnded = true

			if result.IsDraw {
				r.endGameDraw(r.players[0], r.players[1])
			} else if result.IsWhiteWinner && isFirstWhite || !result.IsWhiteWinner && !isFirstWhite {
				r.endGame(r.players[0], r.players[1])
			} else {
				r.endGame(r.players[1], r.players[0])
			}

		}

	}
}

func (r *Room) endGameDraw(player1 *Player, player2 *Player) {
	if !player1.IsClosed() {
		result, _ := MarshalToMessage(ResultMsg, &ResultMessage{"draw", player1.GetScore()})
		errSend := player1.Send(result)
		if errSend != nil {
			singletoneLogger.LogError(errSend)
		}
	}

	if !player2.IsClosed() {
		result, _ := MarshalToMessage(ResultMsg, &ResultMessage{"draw", player2.GetScore()})
		errSend := player2.Send(result)
		if errSend != nil {
			singletoneLogger.LogError(errSend)
		}
	}

	singletoneLogger.LogMessage("Game result: draw")
	r.closeConnections(websocket.CloseNormalClosure, closeNormalGameOverFrame)
	r.isDoneChan <- struct{}{}

}

func (r *Room) endGame(winner *Player, loser *Player) {
	errChangeWinner := winner.playerData.User.AddScore(scoreFactor)
	if errChangeWinner != nil {
		singletoneLogger.LogError(errChangeWinner)
		r.closeConnections(websocket.CloseInternalServerErr, errInternalServerError.Error())
		r.isDoneChan <- struct{}{}
		return
	}

	errChangeLoser := loser.playerData.User.AddScore(-scoreFactor)
	if errChangeLoser != nil {
		singletoneLogger.LogError(errChangeWinner)
		r.closeConnections(websocket.CloseInternalServerErr, errInternalServerError.Error())
		r.isDoneChan <- struct{}{}
		return
	}

	if !winner.IsClosed() {
		result, _ := MarshalToMessage(ResultMsg, &ResultMessage{"win", winner.GetScore() + scoreFactor})
		errSend := winner.Send(result)
		if errSend != nil {
			singletoneLogger.LogError(errSend)
		}
	}

	if !loser.IsClosed() {
		result, _ := MarshalToMessage(ResultMsg, &ResultMessage{"loser", loser.GetScore() - scoreFactor})
		errSend := loser.Send(result)
		if errSend != nil {
			singletoneLogger.LogError(errSend)
		}
	}

	singletoneLogger.LogMessage(fmt.Sprintf("Game Result:\nWINNER - %v,\nLOSER - %v", winner.GetName(), loser.GetName()))
	r.closeConnections(websocket.CloseNormalClosure, closeNormalGameOverFrame)
	r.isDoneChan <- struct{}{}

}

func (r *Room) closeConnections(code int, msg string) {
	for _, p := range r.players {
		p.Close(websocket.CloseNormalClosure, closeNormalGameOverFrame)
	}
}

func (r *Room) listen() {
	closeTimer := time.NewTimer(time.Second * time.Duration(r.closeTime))
	startTimer := time.Now()
LOOP:
	for {
		select {
		case p := <-r.register:

			if !closeTimer.Stop() {
				<-closeTimer.C
			}
			closeTimer.Reset(time.Second * time.Duration(r.closeTime))
			singletoneLogger.LogMessage("Timer was reseted")

			r.players = append(r.players, p)
			singletoneLogger.LogMessage(fmt.Sprintf("RoomID %v: Player was added: %v", r.ID, p.playerData.Name))

			err := p.Send(infoMsgAddedToRoom)
			if websocket.IsUnexpectedCloseError(err) {
				r.isDoneChan <- struct{}{}
			} else if err != nil {
				singletoneLogger.LogError(err)
			}

			if r.isAllPlayersCreated() {
				r.isFullChan <- struct{}{}
			}
		case <-r.isFullChan:
			if !closeTimer.Stop() {
				<-closeTimer.C
			}
			go r.startGame()
			break LOOP
		case <-closeTimer.C:

			for _, p := range r.players {
				p.Close(websocket.CloseTryAgainLater, errCloseWhileSearch.Error())
			}

			singletoneLogger.LogMessage(fmt.Sprintf("Room: %v, was closed due to waiting timeout ", r.ID))
			r.isDoneChan <- struct{}{}
			break LOOP
		}
	}
	<-r.isDoneChan
	finishTime := time.Now().Sub(startTimer)
	singletoneLogger.LogMessage(fmt.Sprintf("%v", finishTime.Seconds()))

	r.game.CloseRoom(r.ID)
}

// color - true - белые, false - черные
func (r *Room) turn(msg *Message, out chan<- *Message, color bool) {
	turn := &TurnMessage{}

	err := msg.UnmarshalData(turn)
	if err != nil {
		log.Println(err)
		out <- errMsgInvalidJSON
	}

	err = r.gameLogic.PlayerTurn(gamelogic.Turn(turn.Turn), color)
	if err != nil {
		errMsg, _ := MarshalToMessage(ErrorMsg, &ErrorMessage{err.Error()})
		out <- errMsg
	} else {
		r.messageToAll(msg)
	}
}
