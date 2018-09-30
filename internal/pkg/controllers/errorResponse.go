package controllers

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"net/http"
)

//easyjson:json
type ErrorResponse struct {
	Error string `json:"error" example:"Some error happened"`
}

func GenerateErrorJSON(errorResponseMessage string) []byte {
	errStruct, err := ErrorResponse{errorResponseMessage}.MarshalJSON()
	if err != nil {
		singletoneLogger.LogError(errors.WithStack(err))
	}
	return errStruct
}

func responseWithError(writer http.ResponseWriter, statusCode int, errorMessage string) {
	writer.WriteHeader(statusCode)
	writer.Write(GenerateErrorJSON(errorMessage))
}
