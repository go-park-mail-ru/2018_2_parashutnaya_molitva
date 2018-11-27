package controllers

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC/core"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/session"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"

	g "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/game"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
)

//go:generate easyjson -pkg

type FindRoom struct {
	Game *g.Game
}

func (fr *FindRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(session.CookieName)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, errNoUser)
	}
	userData, err := fr.Game.GRPCCore.GetUserBySession(r.Context(), &core.Session{Cookie: sessionCookie.Value})
	if err != nil {
		singletoneLogger.LogError(errors.Wrap(err, "can't do grpc request to core for user data by cookie"))
		responseWithError(w, http.StatusUnauthorized, errNoUser)
	}
	b := userData.IsAuth
	if !b {
		responseWithError(w, http.StatusUnauthorized, errNoAuth)
		return
	}

	guid := userData.Guid
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

	roomID, err := fr.Game.InitRoom(guid, params)
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
	sessionCookie, err := r.Cookie(session.CookieName)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, errNoUser)
	}
	userData, err := sg.Game.GRPCCore.GetUserBySession(r.Context(), &core.Session{Cookie: sessionCookie.Value})
	if err != nil {
		singletoneLogger.LogError(errors.Wrap(err, "can't do grpc request to core for user data by cookie"))
		responseWithError(w, http.StatusUnauthorized, errNoUser)
	}
	b := userData.IsAuth
	if !b {
		responseWithError(w, http.StatusUnauthorized, errNoAuth)
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

	sg.Game.InitConnection(userData.Email, userData.Guid, int(userData.Score), conn)
}
