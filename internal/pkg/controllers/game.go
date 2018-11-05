package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/db"
	g "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/game"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"github.com/gorilla/websocket"
)

//go:generate easyjson -pkg

type FindRoom struct {
	Game *g.Game
}

func (gr *FindRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := isAuth(r)
	singletoneLogger.LogMessage("Request")
	if !b {
		responseWithError(w, http.StatusUnauthorized, errNoAuth)
		return
	}

	guid := userGuid(r)
	if guid == "" {
		responseWithError(w, http.StatusUnauthorized, errNoUser)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, errParseRequestBody)
		return
	}
	r.Body.Close()

	params := g.RoomParameters{}
	singletoneLogger.LogMessage(string(body))
	err = params.UnmarshalJSON(body)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, errParseJSON)
		return
	}

	singletoneLogger.LogMessage("here")

	// Пока только один параметр, так что просто кидаем ошибку без названия параметра
	_, err = params.Validate()
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	roomID, err := gr.Game.InitRoom(guid, params)
	if err != nil {
		responseWithError(w, http.StatusConflict, err.Error())
		return
	}
	responseWithOk(w, &FindUserResponse{roomID})
}

//easyjson:json
type FindUserResponse struct {
	RoomID string `json:"roomid" example:"xxx-xx-xxxxx"`
}

type StartGame struct {
	Game     *g.Game
	Upgrader *websocket.Upgrader
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

	// Уже есть CORS middleware, который отклоняет запросы с неразрешенных Origin
	sg.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := sg.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}

	sg.Game.InitConnection(u.Email, u.Score, &u, conn)
}
