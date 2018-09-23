package db

import "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"

const (
	configFilename = "db.json"
)

type dbConfigData struct {
	MongoHost string
	MongoPort string
}

func init() {
	err := jsonConfigReader.Read(configFilename, &dbConfig)
	if err != nil {
		singletoneLogger.LogError(err)
	}
}
