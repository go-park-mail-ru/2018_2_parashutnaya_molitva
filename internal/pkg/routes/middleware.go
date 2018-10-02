package routes

import "net/http"

type Middleware func(h http.Handler) http.Handler

func middlewareChain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
