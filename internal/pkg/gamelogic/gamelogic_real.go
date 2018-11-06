package gamelogic

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/chess"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/randomGenerator"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

type ChessGameLogic struct {
	isWhiteTurn  bool
	gameDuration int

	whiteRemainingTime time.Duration
	blackRemainingTime time.Duration

	whiteTurn chan *Turn
	blackTurn chan *Turn

	turnResponse chan error

	chessEngine ChessEngine
}

func NewChessGameLogic(gameDuration int) *ChessGameLogic {
	return &ChessGameLogic{
		isWhiteTurn:        true,
		gameDuration:       gameDuration,
		whiteRemainingTime: time.Second * time.Duration(gameDuration),
		blackRemainingTime: time.Second * time.Duration(gameDuration),
		whiteTurn:          make(chan *Turn, 1),
		blackTurn:          make(chan *Turn, 1),
		turnResponse:       make(chan error, 1),
		chessEngine:        chess.NewGame(),
	}
}

func (c *ChessGameLogic) start(resultChan chan<- Result) {

	blackTimer := time.NewTimer(c.blackRemainingTime)
	if !blackTimer.Stop() {
		<-blackTimer.C
	}
	whiteTimer := time.NewTimer(c.whiteRemainingTime)

	whitePrevTime := time.Now()
	var whiteAfterStopTime time.Time

	var blackPrevTime time.Time
	var blackAfterStopTime time.Time

	isPrevWhite := false

	for !c.chessEngine.IsGameOver() {
		select {
		case <-whiteTimer.C:

			blackTimer.Stop()

			singletoneLogger.LogMessage("White is timeouted")
			resultChan <- Result{false, false}
			return
		case <-blackTimer.C:

			if !whiteTimer.Stop() {
				<-whiteTimer.C
			}

			singletoneLogger.LogMessage("Black is timeouted")
			resultChan <- Result{true, false}
			return
		case turn := <-c.whiteTurn:
			singletoneLogger.LogMessage("White turn: " + string(*turn))
			err := c.chessEngine.Move(string(*turn))
			if err != nil {
				c.turnResponse <- err
				continue
			}
			if !whiteTimer.Stop() {
				<-whiteTimer.C
			}

			whiteAfterStopTime = time.Now()
			c.whiteRemainingTime = c.whiteRemainingTime - whiteAfterStopTime.Sub(whitePrevTime)

			singletoneLogger.LogMessage(fmt.Sprintf("White remaining time: %v", c.whiteRemainingTime.Seconds()))

			blackTimer.Reset(c.blackRemainingTime)
			blackPrevTime = time.Now()
			c.turnResponse <- nil
			isPrevWhite = true
		case turn := <-c.blackTurn:
			singletoneLogger.LogMessage("Black turn: " + string(*turn))
			err := c.chessEngine.Move(string(*turn))
			if err != nil {
				c.turnResponse <- err
				continue
			}

			if !blackTimer.Stop() {
				<-blackTimer.C
			}

			blackAfterStopTime = time.Now()
			c.blackRemainingTime = c.blackRemainingTime - blackAfterStopTime.Sub(blackPrevTime)

			singletoneLogger.LogMessage(fmt.Sprintf("Black remaining time: %v", c.blackRemainingTime.Seconds()))

			whiteTimer.Reset(c.whiteRemainingTime)
			whitePrevTime = time.Now()
			c.turnResponse <- nil
			isPrevWhite = false
		}
	}

	if c.chessEngine.IsCheckmate() {
		resultChan <- Result{isPrevWhite, false}
	} else if c.chessEngine.IsStalemate() || c.chessEngine.IsInsufficientMaterial() {
		resultChan <- Result{false, true}
	}
}

func (c *ChessGameLogic) PlayerTurn(turn Turn, color bool) error {
	if color != c.isWhiteTurn {
		return errNotYourTurn
	}
	switch color {
	case true:
		c.whiteTurn <- &turn
	case false:
		c.blackTurn <- &turn
	}

	err := <-c.turnResponse
	if err != nil {
		return err
	}
	c.chessEngine.PrintBoard()
	c.chessEngine.PrintLegalMoves()
	c.isWhiteTurn = !c.isWhiteTurn
	return nil
}

func (c *ChessGameLogic) Start() (bool, <-chan Result) {
	resultChan := make(chan Result, 1)
	go c.start(resultChan)

	randomBool := randomGenerator.RandomBool()
	return randomBool, resultChan
}
