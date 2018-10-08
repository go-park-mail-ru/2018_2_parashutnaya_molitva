package auth

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

var (
	authConfig = &config.AuthConfig{}
)

func init() {
	err := config.GetConfig("auth", authConfig)
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
