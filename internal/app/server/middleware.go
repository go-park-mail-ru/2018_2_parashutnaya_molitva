package server

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

type middleware func(h http.HandlerFunc) http.HandlerFunc

func auth() middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			singletoneLogger.LogMessage("Auth middleware")
			ctx := context.WithValue(req.Context(), "isAuth", true)

			h.ServeHTTP(res, req.WithContext(ctx))
		}
	}
}

func middlewareChain(h http.HandlerFunc, middlewares ...middleware) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
