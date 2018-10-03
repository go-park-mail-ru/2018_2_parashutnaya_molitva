package routes

import (
	"strings"

	"github.com/pkg/errors"
)

type ParserResult struct {
	vars map[string]string
}

var (
	errInvalidURL = errors.New("Invalid url")
)
var (
	emptyParserResult = ParserResult{}
)

func parseURL(raw string, pattern string) (bool, ParserResult, error) {
	raw = strings.TrimSuffix(raw, "/")
	pattern = strings.TrimSuffix(pattern, "/")

	if raw != "" && raw[len(raw)-1] == '/' || pattern != "" && pattern[len(pattern)-1] == '/' {
		return false, emptyParserResult, errInvalidURL
	}

	if raw == pattern {
		return true, emptyParserResult, nil
	}

	rawPaths := strings.Split(raw, "/")
	patternPaths := strings.Split(pattern, "/")

	parseResult := ParserResult{
		vars: make(map[string]string),
	}

	for idx, rawPath := range rawPaths {
		if rawPath == patternPaths[idx] {
			continue
		}

		if strings.HasPrefix(patternPaths[idx], ":") && rawPath != "" {
			parseResult.vars[patternPaths[idx][1:]] = rawPath
			continue
		}

		switch patternPaths[idx] {
		case "*":
			return true, emptyParserResult, nil
		default:
			return false, emptyParserResult, nil
		}
	}

	if len(rawPaths) != len(patternPaths) {
		return false, emptyParserResult, nil
	}

	return true, parseResult, nil
}
