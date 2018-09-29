package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/auth"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/session"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
)

func GetSesson(w http.ResponseWriter, r *http.Request) {
	b := r.Context().Value("isAuth")
	if val, _ := b.(bool); val {
		singletoneLogger.LogMessage("isAuth")
	}

	w.Write([]byte("Signin Page"))
}

type signInRequest struct {
}

type SignInResponse struct {
	Error  string               `json:"error"`
	Result SignInResponseResult `json:"result"`
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
// @Success 200
// @Failure 404 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /session [post]
func SignIn(w http.ResponseWriter, r *http.Request) {
	b := r.Context().Value("isAuth").(bool)
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
	w.WriteHeader(http.StatusOK)
}

type SignInParameters struct {
	Email    string `json:"email" example:"test@mail.ru"`
	Password string `json:"password" example:"1234qwerty"`
}
