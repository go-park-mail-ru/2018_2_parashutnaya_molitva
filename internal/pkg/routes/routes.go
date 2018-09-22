package routes

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

type Route struct {
	PathName string
	Handler  http.Handler
	matchers []matcher
}

func (r *Route) addMatcher(m matcher) *Route {
	if m != nil {
		r.matchers = append(r.matchers, m)
	}
	return r
}

func (r *Route) Match(req *http.Request) (bool, error) {
	for _, matcher := range r.matchers {
		if ok, err := matcher.match(req); err != nil {
			return false, err
		} else if ok != true {
			return false, nil
		}
	}
	return true, nil
}

func (r *Route) Path(path string) *Route {
	pathMatcher := &PathMatcher{
		Path: path,
	}
	return r.addMatcher(pathMatcher)
}

func (r *Route) Method(methods ...string) *Route {
	methodMatcher := &MethodMatcher{
		Methods: make([]string, 0, len(methods)),
	}
	for _, method := range methods {
		methodMatcher.Methods = append(methodMatcher.Methods, method)
	}

	return r.addMatcher(methodMatcher)
}

type matcher interface {
	match(req *http.Request) (bool, error)
}

type MethodMatcher struct {
	Methods []string
}

func (m *MethodMatcher) match(req *http.Request) (bool, error) {
	reqMethod := req.Method
	for _, method := range m.Methods {
		if method == reqMethod {
			return true, nil
		}
	}

	return false, fmt.Errorf("%s, doesn't match any of: %v", reqMethod, m.Methods)
}

var (
	errPageNotFound = errors.New("Page not found, 404")
)

type PathMatcher struct {
	Path string
}

func (p *PathMatcher) match(req *http.Request) (bool, error) {
	if req.URL.String() == p.Path {
		return true, nil
	}
	singletoneLogger.LogMessage(req.URL.String() + p.Path)
	singletoneLogger.LogError(errPageNotFound)
	return false, nil
}

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
		if ok, err := route.Match(req); err != nil {
			return false, err
		} else if ok == true {
			r.handler = route.Handler
			return true, nil
		}
	}

	return false, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.routes == nil || r.handler == nil {
		singletoneLogger.LogError(errRouterNotCreated)
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

func (r *Router) HandleFunc(path string, handlerFunc http.HandlerFunc) *Route {
	if r.routes == nil {
		singletoneLogger.LogError(errRouterNotCreated)
		return nil
	}
	route := &Route{
		Handler:  handlerFunc,
		PathName: path,
	}
	r.routes = append(r.routes, route)
	singletoneLogger.LogMessage("Added route: " + route.PathName)
	return route.Path(path)
}
