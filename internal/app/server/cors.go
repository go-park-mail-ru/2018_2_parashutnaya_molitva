package server

import (
	"net/http"
	"strconv"
	"strings"
)

func corsMiddleware(h http.Handler) http.Handler {
	var mw http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		val, ok := req.Header["Origin"]
		if ok {
		LOOP:
			for _, origin := range corsData.AllowOrigins {
				if origin == val[0] {
					res.Header().Set("Access-Control-Allow-Origin", origin)
					res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
					break LOOP
				}
			}
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
