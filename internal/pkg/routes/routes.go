package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
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
	return r.addMatcher(NewPathMatcher(path))
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
	Path           string
	varName        string
	withoutVarName string
}

func parsePathVars(path string) (string, string) {
	s := strings.TrimSuffix(path, "/")
	if s[len(s)-1] == '/' {
		return "", path
	}

	urlPaths := strings.Split(s, "/")
	log.Printf("%#v", urlPaths)
	if urlPaths[len(urlPaths)-1][0] == ':' {
		withoutVarName := strings.Join(urlPaths[:len(urlPaths)-1], "/")
		return urlPaths[len(urlPaths)-1][1:], withoutVarName
	}

	return "", path
}

func NewPathMatcher(path string) *PathMatcher {
	varName, withoutVarName := parsePathVars(path)
	return &PathMatcher{
		Path:           path,
		varName:        varName,
		withoutVarName: withoutVarName,
	}
}

func (p *PathMatcher) parseVarURL(url string) (string, string) {
	if p.varName == "" {
		return "", p.Path
	}
	s := strings.TrimSuffix(url, "/")

	if s[len(s)-1] == '/' {
		return "", p.Path
	}

	urlPaths := strings.Split(s, "/")
	urlVar := urlPaths[len(urlPaths)-1]

	return urlVar, strings.Join(urlPaths[:len(urlPaths)-1], "/")

}

const (
	contextVarKey = iota
)

func (p *PathMatcher) match(req *http.Request) (bool, error) {
	reqURL := req.URL.String()
	if reqURL == p.Path {
		return true, nil
	}

	value, withoutVarName := p.parseVarURL(req.URL.String())
	if value != "" && withoutVarName == p.withoutVarName {

		ctx := context.WithValue(req.Context(), contextVarKey, map[string]string{p.varName: value})
		*req = *req.WithContext(ctx)
		return true, nil
	}

	return false, nil
}

func GetVar(req *http.Request) (map[string]string, bool) {
	value, ok := req.Context().Value(contextVarKey).(map[string]string)
	return value, ok
}
