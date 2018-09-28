package controllers

import (
	"net/http"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

// api/signin
func SignIn(w http.ResponseWriter, r *http.Request) {
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

func User(w http.ResponseWriter, r *http.Request) {
	value, ok := routes.GetVar(r)
	if ok {
		w.Write([]byte(value["id"]))
	}
}
