package routes

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"

	"github.com/pkg/errors"
)

var (
	errRouterNotCreated = errors.New("Router wasn't created")
)

type Router struct {
	routes  []*Route
	handler http.Handler
}

func NewRouter(h http.Handler) *Router {
	return &Router{
		routes:  make([]*Route, 0),
		handler: h,
	}
}

func (r *Router) Match(req *http.Request) (bool, error) {
	for _, route := range r.routes {
		if route.Match(req) {
			r.handler = route.Handler
			return true, nil
		}
	}

	return false, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.routes == nil || r.handler == nil {
		log.Printf(errRouterNotCreated.Error())
		return
	}

	ok, err := r.Match(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if ok != true {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errPageNotFound.Error()))
		return
	}

	r.handler.ServeHTTP(w, req)
}

func (r *Router) HandleFuncWithMiddleware(path string, handlerFunc http.HandlerFunc) *Route {
	handlerFunc = middlewareChain(handlerFunc, authMiddleware)
	return r.HandleFunc(path, handlerFunc)
}

func (r *Router) HandleFunc(path string, handlerFunc http.HandlerFunc) *Route {
	return r.Handle(path, handlerFunc)
}

func (r *Router) Handle(path string, handlerFunc http.Handler) *Route {
	if r.routes == nil {
		singletoneLogger.LogError(errRouterNotCreated)
		return nil
	}
	route := &Route{
		Handler:  handlerFunc,
		PathName: path,
	}
	r.routes = append(r.routes, route)
	return route.path(path)
}
