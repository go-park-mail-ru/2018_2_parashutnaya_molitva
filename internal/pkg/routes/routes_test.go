package routes

import (
	"reflect"
	"testing"
)

type ParseURLTestCase struct {
	url            string
	pattern        string
	isValid        bool
	expectedResult ParserResult
	err            error
}

func TestParseURLRegex(t *testing.T) {
	tests := []*ParseURLTestCase{
		&ParseURLTestCase{
			url:            "/",
			pattern:        "*",
			isValid:        true,
			err:            nil,
			expectedResult: emptyParserResult,
		},
		&ParseURLTestCase{
			url:            "/asdasd/asdsad/asd",
			pattern:        "*",
			isValid:        true,
			err:            nil,
			expectedResult: emptyParserResult,
		},
		&ParseURLTestCase{
			url:            "/asd",
			pattern:        "/asd/*",
			isValid:        false,
			err:            nil,
			expectedResult: emptyParserResult,
		},
		&ParseURLTestCase{
			url:            "/asd/asd/asdasdfasd/",
			pattern:        "/asd/*",
			isValid:        true,
			err:            nil,
			expectedResult: emptyParserResult,
		},
		&ParseURLTestCase{
			url:            "/asd/asd//",
			pattern:        "/asd/*",
			isValid:        false,
			err:            errInvalidURL,
			expectedResult: emptyParserResult,
		},
	}

	for idx, test := range tests {
		isValid, result, err := parseURL(test.url, test.pattern)
		if isValid != test.isValid {
			t.Fatalf("[%v] %#v was expected but got: %v", idx, test.isValid, isValid)
		}
		if test.err == nil && err != nil {
			t.Fatalf("[%v] Unexpected error: %#v", idx, err.Error())
		}

		if test.err != err {
			t.Fatalf("[%v] Error: %#v was expected but got: %v", idx, test.err, err)
		}

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Fatalf("[%v] %#v was expected but got: %v", idx, test.expectedResult, result)
		}
	}
}

func TestParserURLCommon(t *testing.T) {
	tests := []*ParseURLTestCase{
		&ParseURLTestCase{
			url:            "/",
			pattern:        "/",
			isValid:        true,
			err:            nil,
			expectedResult: emptyParserResult,
		},
		&ParseURLTestCase{
			url:            "/asd//asd",
			pattern:        "/asd//asd",
			isValid:        true,
			err:            nil,
			expectedResult: emptyParserResult,
		},
		&ParseURLTestCase{
			url:            "/",
			pattern:        "",
			isValid:        false,
			err:            nil,
			expectedResult: emptyParserResult,
		},
	}

	for idx, test := range tests {
		isValid, result, err := parseURL(test.url, test.pattern)
		if isValid != test.isValid {
			t.Fatalf("[%v] %#v was expected but got: %v", idx, test.isValid, isValid)
		}
		if test.err == nil && err != nil {
			t.Fatalf("[%v] Unexpected error: %#v", idx, err.Error())
		}

		if test.err != err {
			t.Fatalf("[%v] Error: %#v was expected but got: %v", idx, test.err, err)
		}

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Fatalf("[%v] %#v was expected but got: %v", idx, test.expectedResult, result)
		}
	}
}

func TestParseURLVars(t *testing.T) {
	tests := []*ParseURLTestCase{
		&ParseURLTestCase{
			url:     "/url/123/",
			pattern: "/url/:id",
			isValid: true,
			err:     nil,
			expectedResult: ParserResult{
				vars: map[string]string{
					"id": "123",
				},
			},
		},
		&ParseURLTestCase{
			url:     "/url/123/",
			pattern: "/url/:id/",
			isValid: true,
			err:     nil,
			expectedResult: ParserResult{
				vars: map[string]string{
					"id": "123",
				},
			},
		},
		&ParseURLTestCase{
			url:     "/url/123",
			pattern: "/url/:id",
			isValid: true,
			err:     nil,
			expectedResult: ParserResult{
				vars: map[string]string{
					"id": "123",
				},
			},
		},
		&ParseURLTestCase{
			url:     "/url/123/256",
			pattern: "/url/:id/:top",
			isValid: true,
			err:     nil,
			expectedResult: ParserResult{
				vars: map[string]string{
					"id":  "123",
					"top": "256",
				},
			},
		},
		&ParseURLTestCase{
			url:            "/url/123/",
			pattern:        "/url/:id/:top",
			isValid:        false,
			err:            nil,
			expectedResult: emptyParserResult,
		},
		&ParseURLTestCase{
			url:            "/url/123//",
			pattern:        "/url/:id/:top",
			isValid:        false,
			err:            errInvalidURL,
			expectedResult: emptyParserResult,
		},
	}

	for idx, test := range tests {
		isValid, result, err := parseURL(test.url, test.pattern)
		if isValid != test.isValid {
			t.Fatalf("[%v] %#v was expected but got: %v", idx, test.isValid, isValid)
		}
		if test.err == nil && err != nil {
			t.Fatalf("[%v] Unexpected error: %#v", idx, err.Error())
		}

		if test.err != err {
			t.Fatalf("[%v] Error: %#v was expected but got: %v", idx, test.err, err)
		}

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Fatalf("[%v] %#v was expected but got: %v", idx, test.expectedResult, result)
		}
	}
}
