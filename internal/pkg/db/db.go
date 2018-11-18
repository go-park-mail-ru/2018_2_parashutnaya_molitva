package db

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	jsonConfigReader = config.JsonConfigReader{}
	dbConfig         dbConfigData
	instance         *mgo.Database // инстанс синглтона
	once             sync.Once     // Магия для реализации singleton
)

func GetInstance() *mgo.Database {
	once.Do(func() {
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
