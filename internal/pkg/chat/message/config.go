package message

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"time"
)

const (
	configFilename = "chat.json"
)

type chatDatabase struct {
	authConfig        *chatConfigData
	collectionMessage *mgo.Collection
	collectionDialog *mgo.Collection
	session           *mgo.Session
}

var database chatDatabase

type chatConfigData struct {
	MongoHost       string
	MongoPort       string
	MongoUser       string
	MongoPassword   string
}

var (
	jsonConfigReader = config.JsonConfigReader{}
	chatConfig       chatConfigData
)

func init() {
	err := jsonConfigReader.Read(configFilename, &chatConfig)
	if err != nil {
		singletoneLogger.LogError(err)
	}
	info := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s:%s", chatConfig.MongoHost, chatConfig.MongoPort)},
		Timeout:  30 * time.Second,
		Database: "chat",
		Username: chatConfig.MongoUser,
		Password: chatConfig.MongoPassword,
	}
	sess, err := mgo.DialWithInfo(info)
	if err != nil {
		singletoneLogger.LogError(errors.WithStack(err))
	}
	database.collectionMessage = sess.DB("chat").C("message")
	database.collectionDialog = sess.DB("chat").C("dialog")
}
