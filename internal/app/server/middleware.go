package server

import (
	"net/http"
)

type middleware func(h http.HandlerFunc) http.HandlerFunc

func middlewareChain(h http.HandlerFunc, middlewares ...middleware) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func wrapHandlerInMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return middlewareChain(h, corsMiddleware, authMiddleware)
}
