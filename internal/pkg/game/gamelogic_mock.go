package game

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

type GLMock struct{}

func (g *GLMock) FirstPlayerTurn(turn Turn) error {
	singletoneLogger.LogMessage(fmt.Sprintf("First player turn: %v", turn))
	return nil
}

func (g *GLMock) SecondPlayerTurn(turn Turn) error {
	singletoneLogger.LogMessage(fmt.Sprintf("Second player turn. %v", turn))
	return nil
}

func (g *GLMock) Start() (bool, <-chan Result) {

	res := make(chan Result, 1)
	go func() {
		<-time.NewTimer(time.Second * 20).C
		res <- Result{true}
	}()

	return true, res
}
