package auth

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

const (
	configFilename = "auth.json"
)

type authConfigData struct {
	MongoHost       string
	MongoPort       string
	TokenExpireTime int
	TokenLength     int
}

var (
	jsonConfigReader = config.JsonConfigReader{}
	authConfig       authConfigData
)

func init() {
	err := jsonConfigReader.Read(configFilename, &authConfig)
	if err != nil {
		singletoneLogger.LogError(err)
	}
}
