package gamelogic

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

type GLMock struct{}

func (g *GLMock) PlayerTurn(turn Turn, color bool) error {
	singletoneLogger.LogMessage(fmt.Sprintf("Color %v(true - white, false - black)  turn: %v", color, turn))
	return nil
}

func (g *GLMock) Start() (bool, <-chan Result) {

	res := make(chan Result, 1)
	go func() {
		<-time.NewTimer(time.Second * 20).C
		res <- Result{true, false}
	}()

	return true, res
}
