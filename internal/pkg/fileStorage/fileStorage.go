package fileStorage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/randomGenerator"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)


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

func GenerateRandomFileName(ext string) (string, error) {
	fileName, err := randomGenerator.RandomString(10)
	if err != nil {
		singletoneLogger.LogError(err)
		return "", err
	}
	fileName = strings.Join([]string{fileName, ext}, ".")
	return fileName, nil
}
