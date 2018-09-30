package controllers

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/fileStorage"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	userModel "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"net/http"
	"path/filepath"
)

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5 << 20) // 5 MB
	file, headers, err := r.FormFile("file")
	if err != nil {
		singletoneLogger.LogError(err)
	}
	user := userModel.User{Guid: userModel.CreateIdFromString(userGuid(r))}
	fileName, err := fileStorage.GenerateRandomFileName(filepath.Ext(headers.Filename))
	if err != nil {
		singletoneLogger.LogError(err)
	}
	err = fileStorage.UploadFile(file, fileName)
	if err != nil {
		singletoneLogger.LogError(err)
	}
	err = user.ChangeAvatar(fileName)
	if err != nil {
		singletoneLogger.LogError(err)
	}
}
