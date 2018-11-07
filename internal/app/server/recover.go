package server

import (
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

func recoverPanicMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				singletoneLogger.LogMessage("Recover")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}

	return mw
}
