package auth

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Status int

const (
	_ Status = iota
	statusError
	statusNotExist
	statusBadToken
	statusExpired
	statusOk
)

type authDatabase struct {
	authConfig *authConfigData
	collection *mgo.Collection
	session    *mgo.Session
}

var database authDatabase

func init() {
	sess, err := mgo.Dial(fmt.Sprintf("mongodb://%s:%s", authConfig.MongoHost, authConfig.MongoPort))
	if err != nil {
		singletoneLogger.LogError(errors.WithStack(err))
	}
	database.collection = sess.DB("auth").C("sessions")
}

type UserAuthData struct {
	Guid       bson.ObjectId `bson:"_id"` // global user id
	Token      string        `bson:"token"`
	ExpireDate time.Time     `bson:"expire_date"`
}

func create(guid string) (string, error) {
	if !bson.IsObjectIdHex(guid) {
		return "", errors.New("Bad user id")
	}
	token, err := generateToken()
	if err != nil {
		return token, nil
	}
	newItem := &UserAuthData{
		Guid:       bson.ObjectIdHex(guid),
		Token:      token,
		ExpireDate: time.Now().Add(time.Duration(authConfig.TokenExpireTime) * time.Second),
	}
	return token, database.collection.Insert(newItem)

}

func updateOrAddIfNotExist(guid string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return token, nil
	}
	_, err = database.collection.UpsertId(bson.ObjectIdHex(guid),
		bson.M{
			"token":       token,
			"expire_date": time.Now().Add(time.Duration(authConfig.TokenExpireTime) * time.Second),
		})
	return token, err
}

func find(guid string) (*UserAuthData, error) {
	user := &UserAuthData{}
	err := database.collection.Find(bson.M{"_id": bson.ObjectIdHex(guid)}).One(&user)
	return user, err
}

func check(guid string, token string) (Status, error) {
	user, err := find(guid)

	if err != nil && err == mgo.ErrNotFound {
		return statusNotExist, nil
	}
	if err != nil {
		return statusError, err
	}
	if user.Token != token {
		return statusBadToken, nil
	}

	if user.ExpireDate.Before(time.Now()) {
		return statusExpired, nil
	}

	return statusOk, nil
}

func reset(guid string) error {
	token, err := generateToken()
	if err != nil {
		return nil
	}
	_, err = database.collection.UpsertId(bson.ObjectIdHex(guid),
		bson.M{
			"token":       token,
			"expire_date": time.Now(),
		})
	return err
}

func remove(guid string) error {
	return database.collection.RemoveId(bson.ObjectIdHex(guid))
}
