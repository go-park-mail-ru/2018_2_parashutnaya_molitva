package controllers

import (
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
)

func User(w http.ResponseWriter, r *http.Request) {
	value, ok := routes.GetVar(r)
	if ok {
		w.Write([]byte(value["id"]))
		w.Write([]byte(value["kek"]))
	}
}

func Foo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.String()))
}
