package routes

import (
	"testing"
)

type PathVarMatcherTestCase struct {
	path           string
	varName        string
	value          string
	withoutVarName string
	url            string
}

func TestPathVarMatcher(t *testing.T) {
	tests := []*PathVarMatcherTestCase{
		&PathVarMatcherTestCase{
			path:           "/url/:param",
			varName:        "param",
			url:            "/url/1",
			withoutVarName: "/url",
			value:          "1",
		},
		&PathVarMatcherTestCase{
			path:           "/url/",
			varName:        "",
			url:            "/url/1",
			withoutVarName: "/url/",
			value:          "",
		},
	}

	for idx, test := range tests {
		testPathMatcher := pathMatcher{
			Path:    test.path,
			varName: test.varName,
		}
		value, withoutVarName := testPathMatcher.parseVarURL(test.url)

		if value != test.value {
			t.Fatalf("[%v] Value %#v was expected, but get %#v", idx, test.value, value)
		}

		if withoutVarName != test.withoutVarName {
			t.Fatalf("[%v] WithOutVarName %#v was expected, but get %#v", idx, test.withoutVarName, withoutVarName)
		}

	}

}

type parsePathVarsTestCase struct {
	path           string
	varName        string
	withoutVarName string
}

func TestParsePathVars(t *testing.T) {
	tests := []*parsePathVarsTestCase{
		&parsePathVarsTestCase{
			path:           "/url/:id",
			varName:        "id",
			withoutVarName: "/url",
		},
		&parsePathVarsTestCase{
			path:           "/url/",
			varName:        "",
			withoutVarName: "/url/",
		},
	}

	for idx, test := range tests {
		varName, withoutVarName := parsePathVars(test.path)
		if varName != test.varName {
			t.Fatalf("[%v] Varname: %#v  was expected, but get %#v", idx, test.varName, varName)
		}

		if withoutVarName != test.withoutVarName {
			t.Fatalf("[%v] withoutVarName: %#v  was expected, but get %#v", idx, test.withoutVarName, withoutVarName)
		}
	}
}
