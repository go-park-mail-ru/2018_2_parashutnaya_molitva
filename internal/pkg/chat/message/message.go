package message

import (
	"github.com/globalsign/mgo/bson"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"time"
)

//easyjson:json
type Message struct {
	FromGuid string `bson:"from_guid" json:"from_guid"`
	FromLogin string `bson:"from_login" json:"from_login"`
	ToGuid string `bson:"to_guid" json:"to_guid"`
	ToLogin string `bson:"to_login" json:"to_login"`
	Text string `bson:"test" json:"text"`
	Time int64 `bson:"time" json:"time"`
}

type Dialog struct {
	FromLogin string `bson:"from_login" json:"from_login"`
	ToLogin string `bson:"to_login" json:"to_login"`
}

func getUnixMicroseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (message *Message) Save() error {
	if err := message.Validate(); err != nil {
		singletoneLogger.LogError(err)
		return errors.Wrap(err, "Can't save message")
	}
	message.Time = getUnixMicroseconds()
	err := database.collectionMessage.Insert(message)
	if err != nil {
		singletoneLogger.LogError(err)
		return errors.Wrap(err, "Can't save message")
	}
	addDialogIfNotExist(Dialog{message.FromLogin, message.ToLogin})
	return nil
}

func addDialogIfNotExist(dialog Dialog) {
	if dialog.ToLogin == "" ||  dialog.FromLogin == "" {
		return
	}
	dialogReverse := Dialog{dialog.ToLogin, dialog.FromLogin}
	err := database.collectionDialog.Insert(dialog, dialogReverse)
	if err != nil {
		singletoneLogger.LogError(err)
	}
	if err != nil {
		singletoneLogger.LogError(err)
	}
}

func (message *Message) Validate() error {
	if message.Text == "" {
		return errors.New("Empty text")
	}
	if message.FromLogin == "" && message.ToLogin != "" {
		return errors.New("Cant send message from anonymous user to dialog")
	}
	return nil
}

func GetDialogs(fromLogin string) ([]Dialog, error) {
	var result []Dialog
	err := database.collectionDialog.Find(bson.M{"from_login": fromLogin}).All(&result)
	if err != nil {
		singletoneLogger.LogError(err)
		return nil, errors.Wrap(err, "Can't get dialogs")
	}
	return result, nil
}

func GetGlobalMessages(fromTime int64, limit int) ([]Message, error) {
	if limit > 50 || limit == 0 {
		limit = 10
	}
	if fromTime == 0 {
		fromTime = getUnixMicroseconds()
	}
	var result []Message
	err := database.collectionMessage.Find(bson.M{"time": bson.M{"$lte": fromTime}, "to_login": ""}).Sort("-time").Limit(limit).All(&result)
	if err != nil {
		singletoneLogger.LogError(err)
		return nil, errors.Wrap(err, "Can't get global messages")
	}
	return result, nil
}

func GetDialogMessages(fromTime int64, limit int, fromLogin string, toLogin string) ([]Message, error) {
	if limit > 50 || limit == 0 {
		limit = 10
	}
	if fromTime == 0 {
		fromTime = getUnixMicroseconds()
	}
	var result []Message
	authors := []string{fromLogin, toLogin}
	err := database.collectionMessage.Find(bson.M{"time": bson.M{"$lte": fromTime},
		"from_login": bson.M{"$in": authors},
		"to_login": bson.M{"$in": authors}}).
		Sort("-time").Limit(limit).All(&result)
	if err != nil {
		singletoneLogger.LogError(err)
		return nil, errors.Wrap(err, "Can't get dialog messages")
	}
	return result, nil
}