package config

import (
	"os"
	"path/filepath"
)

const configsDir = "configs"

type ConfigReader interface {
	Read(filename string, structure interface{}) error
}

func ProjectPath() string {
	projectPath := "/src/github.com/go-park-mail-ru/2018_2_parashutnaya_molitva"
	fullPath := filepath.Join(os.Getenv("GOPATH"), projectPath)
	return fullPath
}

// формирует полный путь до конфига
func configsPath(filename string) (string, error) {
	fullPath := filepath.Join(ProjectPath(), configsDir, filename)
	return fullPath, nil
}
