package fileStorage

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var StoragePath string

func init() {
	StoragePath = filepath.Join(config.ProjectPath(), "storage")
}

func UploadFile(fileFromRequest multipart.File, fileName string) error {
	fileToSave, err := os.OpenFile(filepath.Join(StoragePath, fileName), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	defer fileToSave.Close()
	_, err = io.Copy(fileToSave, fileFromRequest)
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	return nil
}
