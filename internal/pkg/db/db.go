package db

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"sync"
)

var (
	jsonConfigReader = config.JsonConfigReader{}
	dbConfig         dbConfigData
	instance         *mgo.Database // инстанс синглтона
	once             sync.Once     // Магия для реализации singleton
)

func getInstance() *mgo.Database {
	once.Do(func() {
		sess, err := mgo.Dial(fmt.Sprintf("mongodb://%s:%s", dbConfig.MongoHost, dbConfig.MongoPort))
		if err != nil {
			singletoneLogger.LogError(errors.WithStack(err))
		}
		instance = sess.DB("main")
	})
	return instance
}
