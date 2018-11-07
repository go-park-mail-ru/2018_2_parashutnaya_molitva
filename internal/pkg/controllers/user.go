package controllers

import (
	"io/ioutil"
	"net/http"

	"strconv"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/db"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"github.com/gorilla/mux"
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
	pathVariables := mux.Vars(r)
	if pathVariables == nil {
		responseWithError(w, http.StatusBadRequest, "Bad query")
		return
	}
	guid, ok := pathVariables["guid"]

	if !ok {
		responseWithError(w, http.StatusBadRequest, "Bad query")
		return
	}

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
// @Failure 400 {object} controllers.ErrorFormResponse
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

	field, err := parameters.Validate()
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithFormError(w, http.StatusBadRequest, err.Error(), field)
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

func (c *CreateUserParameters) Validate() (string, error) {
	err := user.ValidateEmail(c.Email)
	if err != nil {
		return "email", err
	}

	err = user.ValidatePassword(c.Password)
	if err != nil {
		return "password", err
	}

	return "", nil
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
// @Failure 400 {object} controllers.ErrorFormResponse
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

	field, err := updateUserData.Validate()
	if err != nil && field != "" {
		responseWithFormError(w, http.StatusBadRequest, err.Error(), field)
		return
	} else if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pathVariables := mux.Vars(r)
	if pathVariables == nil {
		responseWithError(w, http.StatusBadRequest, "Bad query")
		return
	}
	guid, ok := pathVariables["guid"]
	if !ok {
		responseWithError(w, http.StatusBadRequest, "Bad query")
		return
	}

	u, err := user.GetUserByGuid(userGuid(r))
	if u.Guid.Hex() != guid {
		responseWithError(w, http.StatusMethodNotAllowed, "No access rights")
		return
	}
	u.UpdateUser(updateUserData)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithOk(w, responseUserGuidStruct{u.Guid.Hex()})
}

// GetUsersCount godoc
// @Title Count of users
// @Summary Get count of users in a system
// @ID get-user-count
// @Produce  json
// @Success 200 {object} controllers.responseUserGuidStruct
// @Failure 500 {object} controllers.ErrorResponse
// @Router /user/count [get]
func GetUsersCount(w http.ResponseWriter, r *http.Request) {
	count, err := user.GetUsersCount()
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithOk(w, getUsersCountResponse{count})
}

//easyjson:json
type getUsersCountResponse struct {
	Count int `json:"count"`
}

// GetUsersScore godoc
// @Title getScoreOfUsers
// @Summary Returns pairs user email: user score sorted by descendant
// @ID get-user-score
// @Produce  json
// @Param offset query int false "default: 0"
// @Param limit query int false "default: 10"
// @Success 200 {array} controllers.responseUserGuidStruct
// @Failure 500 {object} controllers.ErrorResponse
// @Router /user/score [get]
func GetUsersScore(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	if limit > 100 || limit == 0 {
		limit = 100
	}
	scores, err := user.GetScores(limit, offset)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Unknown error")
		return
	}
	responseWithOk(w, GetUsersScoreResponse{scores})
}

//easyjson:json
type GetUsersScoreResponse struct {
	Scores []user.UserScore `json:"scores"`
}
