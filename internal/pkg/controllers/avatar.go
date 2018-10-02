package controllers

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/fileStorage"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

const (
	MB = 1 << 20
)

var (
	errUploadSize = errors.New("Uploaded file more than 5 mb")
)

// UploadAvatar godoc
// @Title upload avatar
// @Summary upload avatar and returns name of an avatar
// @ID post-avatar
// @Accept  multipart/form-data
// @Produce  json
// @Param avatar formData file true "Avatar file"
// @Success 200 {object} controllers.responseUploadAvatar
// @Failure 400 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /avatar [post]
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 * MB) // 5 MB
	singletoneLogger.LogError(err)

	file, headers, err := r.FormFile("avatar")
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusBadRequest, "Can't parse form ile")
		return
	}

	if headers.Size > 5*MB {
		singletoneLogger.LogError(errUploadSize)
		responseWithError(w, http.StatusBadRequest, errUploadSize.Error())
		return
	}

	fileName, err := fileStorage.GenerateRandomFileName(filepath.Ext(headers.Filename))
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, "Unknown error")
		return
	}
	err = fileStorage.UploadFile(file, fileName)
	if err != nil {
		singletoneLogger.LogError(err)
		responseWithError(w, http.StatusInternalServerError, "Can't upload file")
		return
	}
	responseWithOk(w, responseUploadAvatar{fileName})
}

//easyjson:json
type responseUploadAvatar struct {
	Avatar string `json:"avatar"`
}
