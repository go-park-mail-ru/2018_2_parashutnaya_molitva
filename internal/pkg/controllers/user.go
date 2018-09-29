package controllers

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/db"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"net/http"
)

// GetUser godoc
// @Title Get user
// @Summary Get user data by global user id
// @ID get-user
// @Accept  json
// @Produce  json
// @Param guid query controllers.GetUserParameters true "User id"
// @Success 200 {object} user.User
// @Failure 400 {object} controllers.ErrorResponse
// @Failure 404 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /user/{guid} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	pathVariables, ok := routes.GetVar(r)
	if !ok {
		responseWithError(w, http.StatusBadRequest, "Bad query")
		return
	}
	guid := pathVariables["guid"]
	u, err := user.GetUserByGuid(guid)
	if err != nil && db.IsNotFoundError(err) {
		responseWithError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}
	responseWithOk(w, u)
}

func User2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user1"))
}
