package auth

import (
	"time"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	authConfig *config.AuthConfig
	collection *mgo.Collection
	session    *mgo.Session
}

var database authDatabase

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

func updateOrAddIfNotExist(guid string) (string, time.Time, error) {
	token, err := generateToken()
	if err != nil {
		return token, time.Time{}, nil
	}
	expireTime := time.Now().Add(time.Duration(authConfig.TokenExpireTime) * time.Second)
	_, err = database.collection.UpsertId(bson.ObjectIdHex(guid),
		bson.M{
			"token":       token,
			"expire_date": expireTime,
		})
	return token, expireTime, err
}

func findByGuid(guid string) (*UserAuthData, error) {
	user := &UserAuthData{}
	err := database.collection.Find(bson.M{"_id": bson.ObjectIdHex(guid)}).One(&user)
	return user, err
}

func findByToken(token string) (*UserAuthData, error) {
	user := &UserAuthData{}
	err := database.collection.Find(bson.M{"token": token}).One(&user)
	return user, err
}

func check(token string) (Status, string, error) {
	user, err := findByToken(token)

	if err != nil && err == mgo.ErrNotFound {
		return statusBadToken, "", nil
	}
	if err != nil {
		return statusError, "", err
	}
	if user.Guid == "" {
		return statusNotExist, "", nil
	}

	if user.ExpireDate.Before(time.Now()) {
		return statusExpired, "", nil
	}

	return statusOk, user.Guid.Hex(), nil
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
