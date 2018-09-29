package controllers

import (
	"net/http"
)

func Foo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.String()))
}
