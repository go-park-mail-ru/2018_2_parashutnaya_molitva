package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/db"
	g "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/game"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"github.com/gorilla/websocket"
)

type FindRoom struct {
	Game *g.Game
}

func (gr *FindRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := isAuth(r)
	if !b {
		responseWithError(w, http.StatusUnauthorized, errNoAuth)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, errParseRequestBody)
		return
	}
	defer r.Body.Close()

	params := g.RoomParameters{}
	err = params.UnmarshalJSON(body)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, errParseJSON)
		return
	}

	// Пока только один параметр, так что просто кидаем ошибку без названия параметра
	_, err = params.Validate()
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	roomID := gr.Game.InitRoom(params)
	responseWithOk(w, &FindUserResponse{roomID})
}

//easyjson:json
type FindUserResponse struct {
	RoomID string `json:"roomid" example:"xxx-xx-xxxxx"`
}

func (f *FindUserResponse) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(f)
	return data, err
}

type StartGame struct {
	game     *g.Game
	upgrader *websocket.Upgrader
}

func (sg *StartGame) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := isAuth(r)
	if !b {
		responseWithError(w, http.StatusUnauthorized, errNoAuth)
		return
	}

	guid := userGuid(r)
	u, err := user.GetUserByGuid(guid)
	if err != nil && db.IsNotFoundError(err) {
		responseWithError(w, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	conn, err := sg.upgrader.Upgrade(w, r, nil)
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}

	sg.game.InitConnection(u.Email, u.Score, conn)
}
