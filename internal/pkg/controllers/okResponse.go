package controllers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"net/http"
)

func GenerateMessageJSON(responseMessage json.Marshaler) []byte {
	responseMessageJSON, err := responseMessage.MarshalJSON()
	if err != nil {
		singletoneLogger.LogError(errors.WithStack(err))
	}
	return responseMessageJSON
}

func responseWithOk(writer http.ResponseWriter, responseMessage json.Marshaler) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("TEST", "100500")
	writer.Write(GenerateMessageJSON(responseMessage))
}
