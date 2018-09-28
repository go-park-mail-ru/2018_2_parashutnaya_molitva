package controllers

import "net/http"

func isAuth(r *http.Request) bool {
	return r.Context().Value("isAuth").(bool)
}

func userGuid(r *http.Request) string {
	return r.Context().Value("userGuid").(string)
}
