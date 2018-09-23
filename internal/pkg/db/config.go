package db

const (
	configFilename = "db.json"
)

type dbConfigData struct {
	MongoHost string
	MongoPort string
}
