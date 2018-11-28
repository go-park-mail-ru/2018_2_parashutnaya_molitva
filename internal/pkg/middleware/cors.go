package middleware

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
)

var (
	configReader = config.JsonConfigReader{}
	confifFile   = "cors.json"
	corsData     CorsData
)
func init() {
	err := configReader.Read(confifFile, &corsData)
	if err != nil {
		singletoneLogger.LogError(err)
	}
}

//easyjson:json
type CorsData struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	MaxAge           int
	AllowCredentials bool
}

func CorsMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		val, ok := req.Header["Origin"]
		if ok {
				res.Header().Set("Access-Control-Allow-Origin", val[0])
				res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
		}

		if req.Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
			res.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
			res.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsData.MaxAge))
			return
		}

		h.ServeHTTP(res, req)
	}

	return mw
}
