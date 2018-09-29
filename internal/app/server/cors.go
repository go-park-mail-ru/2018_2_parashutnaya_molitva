package server

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
)

var (
	configReader = config.JsonConfigReader{}
	confifFile   = "cors.json"
	corsData     CorsData
)

type CorsData struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	MaxAge           int
	AllowCredentials bool
}

func corsMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		val, ok := req.Header["Origin"]
		if ok {
			singletoneLogger.LogMessage(fmt.Sprintf("%#v", val))
		}

		h.ServeHTTP(res, req)
	}
}
