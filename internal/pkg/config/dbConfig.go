package config

const dbConfigName = "db"

type DbConfig struct {
	MongoHost     string
	MongoPort     string
	MongoUser     string
	MongoPassword string
}
