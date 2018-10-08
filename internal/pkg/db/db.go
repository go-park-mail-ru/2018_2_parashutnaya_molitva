package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

var (
	instance *mgo.Database // инстанс синглтона
	once     sync.Once     // Магия для реализации singleton
)

func GetInstance() *mgo.Database {
	once.Do(func() {
		dbConfig := config.DbConfig{}
		err := config.GetConfig("db", &dbConfig)
		if err != nil {
			singletoneLogger.LogError(err)
		}

		info := &mgo.DialInfo{
			Addrs:    []string{fmt.Sprintf("%s:%s", dbConfig.MongoHost, dbConfig.MongoPort)},
			Timeout:  30 * time.Second,
			Database: "main",
			Username: dbConfig.MongoUser,
			Password: dbConfig.MongoPassword,
		}
		sess, err := mgo.DialWithInfo(info)
		if err != nil {
			singletoneLogger.LogError(errors.WithStack(err))
		}
		instance = sess.DB("main")
	})
	return instance
}

func IsNotFoundError(err error) bool {
	if (err != nil) && (err.Error() == "not found") {
		return true
	}
	return true
}
