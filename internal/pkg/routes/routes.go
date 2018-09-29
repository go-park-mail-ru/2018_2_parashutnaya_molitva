package routes

import (
	"context"

	"net/http"
	"reflect"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"

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

func (r *Route) Match(req *http.Request) bool {
	for _, matcher := range r.matchers {
		if !matcher.match(req) {
			return false
		}
	}
	return true
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
	match(req *http.Request) bool
}

type MethodMatcher struct {
	Methods []string
}

func (m *MethodMatcher) match(req *http.Request) bool {
	reqMethod := req.Method
	for _, method := range m.Methods {
		if method == reqMethod {
			return true
		}
	}

	return false
}

var (
	errPageNotFound = errors.New("Page not found, 404")
)

type pathMatcher struct {
	Pattern string
}

const (
	contextVarKey = iota
)

func (p *pathMatcher) match(req *http.Request) bool {
	reqURL := req.URL.String()

	isValid, resultParse, err := parseURL(reqURL, p.Pattern)
	if err != nil {
		singletoneLogger.LogError(err)
		return false
	}
	if isValid {
		if !reflect.DeepEqual(resultParse, emptyParserResult) {
			ctx := context.WithValue(req.Context(), contextVarKey, resultParse.vars)
			*req = *req.WithContext(ctx)
		}
		return true
	}

	return false
}

func GetVar(req *http.Request) (map[string]string, bool) {
	value, ok := req.Context().Value(contextVarKey).(map[string]string)
	return value, ok
}
