package gRPC

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

const (
	configFilename = "grpc.json"
)

//easyjson:json
type configData struct {
	CorePort string
	AuthPort string
}

var (
	jsonConfigReader         = config.JsonConfigReader{}
	grpcConfig       configData
)

func init() {
	err := jsonConfigReader.Read(configFilename, &grpcConfig)
	if err != nil {
		singletoneLogger.LogError(err)
	}
}

func GetCorePort() string {
	return grpcConfig.CorePort
}

func GetAuthPort() string {
	return grpcConfig.AuthPort
}