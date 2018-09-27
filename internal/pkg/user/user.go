package user

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/fileStorage"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/randomGenerator"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"gopkg.in/mgo.v2/bson"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type User struct {
	Guid         bson.ObjectId `bson:"_id"`
	Email        string        `bson:"email"`
	Password     string
	HashPassword string `bson:"password"`
	Avatar       string `bson:"avatar"`
	Score        uint   `bson:"score"`
}

func (u *User) UploadAvatar(f multipart.File, fileName string) error {
	ext := filepath.Ext(fileName)
	avatarName, err := randomGenerator.RandomString(10)
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	avatarName = strings.Join([]string{avatarName, ext}, ".")
	err = fileStorage.UploadFile(f, avatarName)
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	u.Avatar = avatarName
	u.updateUser()
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	return nil
}

func LoginUser(email string, password string) (*User, error) {
	u, err := getUserByEmail(email)
	if err != nil {
		singletoneLogger.LogError(err)
		return &User{}, err
	}
	err = checkPasswordByHash(password, u.HashPassword)
	if err != nil {
		singletoneLogger.LogError(err)
		return &User{}, err
	}
	return u, err
}

func (u *User) ChangeEmail(email string) error {
	u.Email = email
	err := collection.UpdateId(bson.M{"_id": u.Guid}, u)
	return err
}

func (u *User) ChangePassword(password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	u.HashPassword = hashedPassword
	err = collection.UpdateId(bson.M{"_id": u.Guid}, u)
	return err
}

func (u *User) AddScore(score uint) error {
	u.Score += score
	err := collection.UpdateId(bson.M{"_id": u.Guid}, u)
	return err
}

func createUser(email string, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		singletoneLogger.LogError(err)
		return &User{}, err
	}
	u := User{
		Guid:         bson.NewObjectId(),
		Email:        email,
		HashPassword: hashedPassword,
	}
	err = collection.Insert(u)
	if err != nil {
		singletoneLogger.LogError(err)
		return &User{}, err
	}
	return u, nil
}

func CreateIdFromString(str string) bson.ObjectId {
	return bson.ObjectIdHex(str)
}
