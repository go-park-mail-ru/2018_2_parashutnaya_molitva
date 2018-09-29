package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
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

func (r *Route) path(pattern string) *Route {
	return r.addMatcher(&pathMatcher{Pattern: pattern})
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

type pathMatcher struct {
	Pattern string
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

func NewPathMatcher(pattern string) *pathMatcher {
	return &pathMatcher{
		Pattern: pattern,
	}
}

const (
	contextVarKey = iota
)

func (p *pathMatcher) match(req *http.Request) (bool, error) {
	reqURL := req.URL.String()

	isValid, resultParse, err := parseURL(reqURL, p.Pattern)
	if err != nil {
		return false, err
	}
	if isValid {
		if !reflect.DeepEqual(resultParse, emptyParserResult) {
			ctx := context.WithValue(req.Context(), contextVarKey, resultParse.vars)
			*req = *req.WithContext(ctx)
		}
		return true, nil
	}

	return false, nil
}

func GetVar(req *http.Request) (map[string]string, bool) {
	value, ok := req.Context().Value(contextVarKey).(map[string]string)
	return value, ok
}
