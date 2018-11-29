package gamelogic

import (
	"github.com/pkg/errors"
	"time"
)

var (
	errNotYourTurn = errors.New("It's not you turn")
)

type Result struct {
	IsWhiteWinner bool
	IsDraw        bool
}

type Turn string

type GameLogic interface {
	PlayerTurn(Turn, bool) (time.Duration, time.Duration, error)  // Turn - ход. 2 параметр, если true - это ход белых, false - это ход черных
	Start() (bool, <-chan Result) // Возвращает рандомный true/false. Можно использовать для того, чтобы определить кто White, а кто Black
	Stop()
}
