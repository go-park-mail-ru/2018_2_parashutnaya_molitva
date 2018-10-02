package controllers

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/auth"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/session"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
)

func Session(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		signIn(w, r)
	case "GET":
		getSession(w, r)
	}
}

// DeleteSession godoc
// @Title Delete session
// @Summary Delete current user session from server
// @ID delete-session
// @Produce  json
// @Success 200 {object} controllers.responseUserGuidStruct
// @Failure 401 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /session [delete]
func DeleteSession(w http.ResponseWriter, r *http.Request) {
	b := isAuth(r)
	if !b {
		responseWithError(w, http.StatusUnauthorized, "Does not authorised")
		return
	}
	cookie, noCookie := r.Cookie(session.CookieName)
	if noCookie != nil {
		responseWithError(w, http.StatusUnauthorized, "Does not authorised")
		return
	}

	guid := userGuid(r)
	err := auth.DeleteSession(guid)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
	}

	cookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

// GetSession godoc
// @Title Get session
// @Summary Get current user of session
// @ID get-session
// @Produce  json
// @Success 200  {object} controllers.responseUserGuidStruct
// @Failure 401 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /session [GET]
func getSession(w http.ResponseWriter, r *http.Request) {
	b := isAuth(r)
	if !b {
		responseWithError(w, http.StatusUnauthorized, "Does not authorised")
		return
	}
	guid := userGuid(r)
	if guid == "" {
		responseWithError(w, http.StatusInternalServerError, "Can't find user")
		return
	}
	responseWithOk(w, responseUserGuidStruct{guid})
}

type SignInResponseResult struct {
}

// SignIn godoc
// @Title Sign in
// @Summary Sign in with your account with email and password, set session cookie
// @ID post-session
// @Accept  json
// @Produce  json
// @Param AuthData body controllers.SignInParameters true "User auth data"
// @Success 200 {object} controllers.responseUserGuidStruct
// @Failure 404 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /session [post]
func signIn(w http.ResponseWriter, r *http.Request) {
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
	parameters := SignInParameters{}
	err = parameters.UnmarshalJSON(body)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, "Can't parse json")
		return
	}
	singletoneLogger.LogMessage(parameters.Email)
	singletoneLogger.LogMessage(parameters.Password)
	u, err := user.LoginUser(parameters.Email, parameters.Password)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	token, timeExpire, err := auth.SetSession(u.Guid.Hex())
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, "Couldn't set the sesson")
		return
	}
	http.SetCookie(w, session.CreateAuthCookie(token, timeExpire))
	responseWithOk(w, responseUserGuidStruct{u.Guid.Hex()})
}

//easyjson:json
type SignInParameters struct {
	Email    string `json:"email" example:"test@mail.ru"`
	Password string `json:"password" example:"1234qwerty"`
}
