package game_test

import (
	"net/http"
	"testing"
	"time"

	g "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/game"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type TestInitConnHandler struct {
	game *g.Game
}

func (th *TestInitConnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	err = th.game.InitConnection("name", 0, c)
	err = c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	if err != nil {
		return
	}
}

type TestCaseInitConn struct {
	params TestCaseInitRoom
	name   string
}

func TestNoErrConnection(t *testing.T) {
	params := []TestCaseInitRoom{
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
	}
	game := g.NewGame()
	fId := game.InitRoom(params[0].params)

	upgrader := websocket.Upgrader{}
	upgrader.Upgrade()
}
