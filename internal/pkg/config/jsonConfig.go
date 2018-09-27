package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type JsonConfigReader struct{}

func (JsonConfigReader) Read(filename string, structure interface{}) error {
	fullFileName, err := configsPath(filename)
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

	err = json.Unmarshal(body, structure)
	if err != nil {
		return err
	}

	return nil
}
