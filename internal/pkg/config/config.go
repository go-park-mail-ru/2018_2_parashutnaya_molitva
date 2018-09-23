package config

import (
	"path/filepath"
)

const configsDir = "configs"

type ConfigReader interface {
	Read(filename string, structure interface{}) error
}

// формирует полный путь до конфига
func normalizeFilepath(filename string) (string, error) {
	binPath, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(binPath, configsDir, filename)
	return fullPath, nil
}
