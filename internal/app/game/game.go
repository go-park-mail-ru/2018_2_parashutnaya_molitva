package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	g "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/game"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/session"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	uuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
	Game *g.Game
}

var users = map[string]string{}
var upgrader = websocket.Upgrader{}

func (ws *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	singletoneLogger.LogMessage("Websocket")
	cookie, noCookie := r.Cookie("sessionid")
	if noCookie != nil {
		singletoneLogger.LogMessage("No cookie")
		return
	}

	username := users[cookie.Value]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}

	err = ws.Game.InitConnection(username, 0, conn)
	if err != nil {
		singletoneLogger.LogError(err)
		errW := conn.WriteJSON(&ErrResponse{err.Error()})
		if errW != nil {
			singletoneLogger.LogError(errW)
		}
	}
}

type GSStruct struct {
	Username     string `json:"username"`
	RoomDuration int    `json:"duration"`
}

type Response struct {
	RoomID string `json:"roomid"`
}

type ErrResponse struct {
	err string `json:"error"`
}

type GameStart struct {
	Game *g.Game
}

func (gs *GameStart) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	singletoneLogger.LogMessage("Request")
	if r.Method != "POST" {
		w.Write([]byte("Invalid method"))
		w.WriteHeader(400)
		return
	}

	singletoneLogger.LogMessage("Request")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		singletoneLogger.LogError(err)
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	gss := &GSStruct{}
	err = json.Unmarshal(body, gss)
	if err != nil {
		singletoneLogger.LogError(err)
		w.WriteHeader(500)
		return
	}

	cookie := uuid.New().String()
	users[cookie] = gss.Username

	httpCookie := session.CreateAuthCookie(cookie, time.Now().Add(time.Duration(10)*time.Minute))
	http.SetCookie(w, httpCookie)

	singletoneLogger.LogMessage(fmt.Sprintf("GSS: %#v", gss))
	roomId := gs.Game.InitRoom(g.RoomParameters{time.Second * time.Duration(gss.RoomDuration)})

	resp, err := json.Marshal(Response{roomId})
	if err != nil {
		singletoneLogger.LogError(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "null")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Write(resp)
	w.WriteHeader(200)
	return
}

func StartGameServer(port int) error {
	stringPort := ":" + strconv.Itoa(port)

	game := g.NewGame()
	webSocket := &WebSocket{game}
	gameStart := &GameStart{game}
	http.Handle("/game", gameStart)
	http.Handle("/game/ws", webSocket)

	singletoneLogger.LogMessage("Server starting at " + stringPort)
	return http.ListenAndServe(stringPort, nil)
}
