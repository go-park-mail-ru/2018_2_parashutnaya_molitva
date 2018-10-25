package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

type GLMock struct{}

func (g *GLMock) FirstPlayerTurn(turn Turn) error {
	singletoneLogger.LogMessage(fmt.Sprintf("First player turn. From: %v, To: %v", turn.from, turn.to))
	return nil
}

func (g *GLMock) SecondPlayerTurn(turn Turn) error {
	singletoneLogger.LogMessage(fmt.Sprintf("Second player turn. From: %v, To: %v", turn.from, turn.to))
	return nil
}
