package game_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
		fmt.Println(err)
		return
	}

	th.game.InitConnection("name", 0, c)
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

	s := httptest.NewServer(http.Handler(&TestInitConnHandler{game}))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	data, _ := json.Marshal(&g.InitMessage{fId})

	if err := ws.WriteJSON(&g.Message{"init", data}); err != nil {
		t.Fatalf("%v", err)
	}
	_, resp, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	resResp, err := g.UnmarshalToMessage(resp)
	if err != nil {
		t.Fatalf("%v", err)
	}
	data, _ = json.Marshal(&g.InfoMessage{"Added to room"})
	expectedMsg := g.NewMessage("info", data)

	if !reflect.DeepEqual(resResp, expectedMsg) {
		t.Fatalf("Expected: %#v, but get %#v\n", expectedMsg, resResp)
	}

	t.Logf("Resp: %#v\n", resp)

}
