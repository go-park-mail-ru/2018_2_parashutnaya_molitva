package controllers

import "net/http"

// api/signup
func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signup Page"))
}
