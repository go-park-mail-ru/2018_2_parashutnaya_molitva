package controllers

import (
	"net/http"
)

// api/signin
func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signin Page"))
}

// api/signup
func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signup Page"))
}
