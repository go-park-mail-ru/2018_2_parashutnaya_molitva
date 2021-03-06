package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

func RecoverPanicMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				singletoneLogger.LogMessage(fmt.Sprintf("Recover from error: %v", r))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}

	return mw
}
