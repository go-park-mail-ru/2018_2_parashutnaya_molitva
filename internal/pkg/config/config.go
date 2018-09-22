package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config interface {
	Read(filename string, structure interface{}) error
}

const configsDir = "configs"

var jsonConfig = new(JsonConfig)

// формирует полный путь до конфига
func normalizeFilepath(filename string) (string, error) {
	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	arrPath := strings.Split(binPath, string(filepath.Separator))
	arrPath = append([]string{"/"}, arrPath...)
	arrPath = append(arrPath[:len(arrPath)-2], configsDir)
	arrPath = append(arrPath, filename)
	fullPath := filepath.Join(arrPath...)
	return fullPath, nil
}

type JsonConfig struct{}

func (JsonConfig) Read(filename string, structure interface{}) error {
	fullFileName, err := normalizeFilepath(filename)
	if err != nil {
		return err
	}
	f, err := os.Open(fullFileName)
	if err != nil {
		return err
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	log.Println(string(body))

	err = json.Unmarshal(body, structure)
	if err != nil {
		return err
	}

	return nil
}
