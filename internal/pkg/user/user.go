package user

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Guid         bson.ObjectId `bson:"_id"`
	Email        string        `bson:"email"`
	Password     string	`bson:"-"`
	HashPassword string `bson:"password"`
	Avatar       string `bson:"avatar"`
	Score        uint   `bson:"score"`
}

func (u *User) ChangeAvatar(avatarName string) error {
	u.Avatar = avatarName
	err := collection.UpdateId(bson.M{"_id": u.Guid}, u)
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	return nil
}

func LoginUser(email string, password string) (User, error) {
	u, err := GetUserByEmail(email)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, err
	}
	err = checkPasswordByHash(password, u.HashPassword)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, err
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

func CreateUser(email string, password string) (User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, err
	}
	u := User{
		Guid:         bson.NewObjectId(),
		Email:        email,
		HashPassword: hashedPassword,
	}
	err = collection.Insert(u)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, err
	}
	return u, nil
}

func (u *User) DeleteUser() error {
	err := collection.RemoveId(u.Guid)
	return err
}

func GetUserByGuid(guid string) (User, error) {
	user := User{}
	err := collection.FindId(bson.M{"_id": CreateIdFromString(guid)}).One(&user)
	return user, err
}

func GetUserByEmail(email string) (User, error) {
	user := User{}
	err := collection.Find(bson.M{"email": email}).One(&user)
	return user, err
}

func CreateIdFromString(str string) bson.ObjectId {
	return bson.ObjectIdHex(str)
}
