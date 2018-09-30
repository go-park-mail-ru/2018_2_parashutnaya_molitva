package server

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/auth"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/session"
)

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
