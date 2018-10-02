package controllers

import (
	"net/http"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
)

//easyjson:json
type ErrorResponse struct {
	Error string `json:"error" example:"Some error happened"`
}

func GenerateErrorJSON(errorResponseMessage string) []byte {
	errStruct, err := ErrorResponse{errorResponseMessage}.MarshalJSON()
	if err != nil {
		singletoneLogger.LogError(errors.WithStack(err))
		return []byte("")
	}
	return errStruct
}

func responseWithError(writer http.ResponseWriter, statusCode int, errorMessage string) {
	writer.WriteHeader(statusCode)
	writer.Write(GenerateErrorJSON(errorMessage))
}

//easyjson:json
type ErrorFormResponse struct {
	Error string `json:"error" example:"Some error happened"`
	Field string `json:"field" example:"email"`
}

func GenerateErrorFormJSON(errorMessage, field string) []byte {
	errStruct, err := ErrorFormResponse{
		Error: errorMessage,
		Field: field,
	}.MarshalJSON()

	if err != nil {
		singletoneLogger.LogError(errors.WithStack(err))
		return []byte("")
	}

	return errStruct
}

func responseWithFormError(writer http.ResponseWriter, statusCode int, errorMessage string, field string) {
	writer.WriteHeader(statusCode)
	writer.Write(GenerateErrorFormJSON(errorMessage, field))
}
