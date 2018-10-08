package config

const authConfigName = "auth"

type AuthConfig struct {
	MongoHost       string
	MongoPort       string
	TokenExpireTime int
	TokenLength     int
	MongoUser       string
	MongoPassword   string
}
