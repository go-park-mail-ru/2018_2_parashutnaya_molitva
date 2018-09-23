package config

import (
	"path/filepath"
	"os"
)

const configsDir = "configs"

type ConfigReader interface {
	Read(filename string, structure interface{}) error
}

// формирует полный путь до конфига
func normalizeFilepath(filename string) (string, error) {
	projectPath := "/src/github.com/go-park-mail-ru/2018_2_parashutnaya_molitva"
	fullPath := filepath.Join(os.Getenv("GOPATH"),projectPath, configsDir, filename)
	return fullPath, nil
}
