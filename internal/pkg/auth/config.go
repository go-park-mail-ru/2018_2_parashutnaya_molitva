package auth

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	configFilename = "auth.json"
)

type authConfigData struct {
	MongoHost       string
	MongoPort       string
	TokenExpireTime int
	TokenLength     int
	MongoUser	string
	MongoPassword string
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
	info := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s:%s", authConfig.MongoHost, authConfig.MongoPort)},
		Timeout:  30 * time.Second,
		Database: "auth",
		Username: authConfig.MongoUser,
		Password: authConfig.MongoPassword,
	}
	sess, err := mgo.DialWithInfo(info)
	if err != nil {
		singletoneLogger.LogError(errors.WithStack(err))
	}
	database.collection = sess.DB("auth").C("sessions")
}
