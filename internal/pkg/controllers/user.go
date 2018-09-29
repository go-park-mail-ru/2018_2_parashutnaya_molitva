package controllers

import (
	"net/http"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
)

func User (w http.ResponseWriter, r *http.Request) {
	pathVariables, ok := routes.GetVar(r)
	if ok {
		w.Write([]byte("user1"))
		w.Write([]byte(pathVariables["id"]))
	}
	return
	switch r.Method {
	case "POST":
		signIn(w,r)
	case "GET":
		getSesson(w,r)

	}
}

func User2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user1"))
}
