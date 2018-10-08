package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var (
	logMode      = true
	configIsRead = false
	once         sync.Once
)

const (
	configsDir     = "configs"
	configFileName = "config.json"
)

func logWrapper(str string) {
	if logMode {
		log.Println(str)
	}
}

func SetDebugLog(debugLog bool) {
	logMode = debugLog
}

func convertMapToConfigStruct(confRaw map[string]interface{}, configStruct interface{}) error {
	return mapstructure.Decode(confRaw, configStruct)
}

func readConfig(filename string) error {
	viper.SetConfigFile(filename)
	return viper.ReadInConfig()
}

func ProjectPath() string {
	projectPath := "/src/github.com/go-park-mail-ru/2018_2_parashutnaya_molitva"
	fullPath := filepath.Join(os.Getenv("GOPATH"), projectPath)
	return fullPath
}

func GetConfig(name string, configStruct interface{}) error {
	once.Do(func() {
		path, err := configsPath(configFileName)
		if err != nil {
			log.Println(err)
		}
		err = readConfig(path)
		if err != nil {
			log.Println(err)
		}
	})
	conf := viper.GetStringMap(name)
	return convertMapToConfigStruct(conf, configStruct)
}

// формирует полный путь до конфига
func configsPath(filename string) (string, error) {
	fullPath := filepath.Join(ProjectPath(), configsDir, filename)
	return fullPath, nil
}
