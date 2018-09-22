package auth

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

const (
	configFilename = "auth.json"
)

var (
	jsonConfigReader = config.JsonConfigReader{}
)

type authConfigData struct {
	MongoHost string
	MongoPort string
}

func init() {
	var authConfig authConfigData
	err := jsonConfigReader.Read(configFilename, &authConfig)
	if err != nil {
		singletoneLogger.LogError(err)
	}
}
