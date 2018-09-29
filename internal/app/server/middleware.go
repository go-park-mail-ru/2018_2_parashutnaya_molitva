package server

import (
	"context"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/auth"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/session"
	"net/http"
)

type middleware func(h http.HandlerFunc) http.HandlerFunc

func authMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		sessionCookie, errNoCookie := req.Cookie(session.CookieName)
		if errNoCookie != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			ctx = context.WithValue(ctx, "userGuid", "")
		} else {
			isAuth, guid, _ := auth.CheckSession(sessionCookie.Value)
			ctx = context.WithValue(ctx, "isAuth", isAuth)
			ctx = context.WithValue(ctx, "userGuid", guid)
		}
		h.ServeHTTP(res, req.WithContext(ctx))
	}
}
func middlewareChain(h http.HandlerFunc, middlewares ...middleware) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
