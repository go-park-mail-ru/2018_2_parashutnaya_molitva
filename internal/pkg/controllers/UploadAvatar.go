package controllers

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	userModel "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"net/http"
)

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5 << 20)
	file, headers, err := r.FormFile("file")

	if err != nil {
		singletoneLogger.LogError(err)
	}
	user := userModel.User{Guid: userModel.CreateIdFromString(userGuid(r))}
	err := user.UploadAvatar(file, headers.Filename)
	if err != nil {

	}
}
