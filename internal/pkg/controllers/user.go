package controllers

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/db"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"net/http"
	"io/ioutil"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

//easyjson:json
type responseUserGuidStruct struct {
	UserGuid string `json:"user_guid"`
}

// GetUser godoc
// @Title Get user
// @Summary Get user data by global user id
// @ID get-user
// @Accept  json
// @Produce  json
// @Param guid query string true "User id"
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
	return
}

// CreateUser godoc
// @Title Create user
// @Summary Sign up with email and password
// @ID post-user
// @Accept  json
// @Produce  json
// @Param data body controllers.CreateUserParameters true "User id"
// @Success 200 {object} controllers.responseUserGuidStruct
// @Failure 400 {object} controllers.ErrorResponse
// @Failure 404 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	b := isAuth(r)
	if b {
		responseWithError(w, http.StatusBadRequest, "Already signed in")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, "Can't parse request body")
		return
	}
	defer r.Body.Close()
	parameters := CreateUserParameters{}
	err = parameters.UnmarshalJSON(body)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, "Can't parse json")
		return
	}
	u, err := user.CreateUser(parameters.Email, parameters.Password)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	responseWithOk(w, responseUserGuidStruct{u.Guid.Hex()})
}

//easyjson:json
type CreateUserParameters struct {
	Email    string `json:"email" example:"test@mail.ru"`
	Password string `json:"password" example:"1234qwerty"`
}

// UpdateUser godoc
// @Title Update user
// @Summary Update current user with data
// @ID put-user
// @Accept  json
// @Produce  json
// @Param data body user.UpdateUserStruct true "updating data"
// @Success 200 {object} controllers.responseUserGuidStruct
// @Failure 400 {object} controllers.ErrorResponse
// @Failure 404 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /user/{guid} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	b := isAuth(r)
	if !b {
		responseWithError(w, http.StatusBadRequest, "Is not authorised")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, "Can't parse request body")
		return
	}
	defer r.Body.Close()
	updateUserData := user.UpdateUserStruct{}
	err = updateUserData.UnmarshalJSON(body)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, "Can't parse json")
		return
	}
	pathVariables, ok := routes.GetVar(r)
	if !ok {
		responseWithError(w, http.StatusBadRequest, "Bad query")
		return
	}
	guid := pathVariables["guid"]
	u, err := user.GetUserByGuid(userGuid(r))
	if u.Guid.Hex() != guid {
		responseWithError(w, http.StatusMethodNotAllowed, "No access rights")
		return
	}

	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	responseWithOk(w, responseUserGuidStruct{u.Guid.Hex()})
}

// GetUsersCount godoc
// @Title Count of users
// @Summary Get count of users in a system
// @ID get-user-count
// @Accept  json
// @Produce  json
// @Success 200 {object} controllers.responseUserGuidStruct
// @Failure 500 {object} controllers.ErrorResponse
// @Router /user/count/ [get]
func GetUsersCount (w http.ResponseWriter, r *http.Request) {
	count, err := user.GetUsersCount()
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithOk(w, getUsersCountResponse{count})
}

//easyjson:json
type getUsersCountResponse struct{
	Count int `json:"count"`
}