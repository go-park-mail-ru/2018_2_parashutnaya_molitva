package user

import (
	simpleErrors "errors"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

//easyjson:json
type User struct {
	Guid         bson.ObjectId `bson:"_id" json:"guid"`
	Email        string        `bson:"email" json:"email"`
	Password     string        `bson:"-" json:"-"`
	HashPassword string        `bson:"password" json:"-"`
	Avatar       string        `bson:"avatar" json:"avatar"`
	Score        int           `bson:"score" json:"score"`
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
	if (err != nil) && (err.Error() == "not found") {
		return User{}, simpleErrors.New("User not found")
	}
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, simpleErrors.New("Internal error")
	}
	err = checkPasswordByHash(password, u.HashPassword)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, simpleErrors.New("Wrong password")
	}
	return u, nil
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

func (u *User) AddScore(score int) error {
	u.Score += score
	err := collection.UpdateId(bson.M{"_id": u.Guid}, u)
	return err
}

func CreateUser(email string, password string) (User, error) {
	isExisting, err := IsUserExisting(email)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, err
	}
	if isExisting {
		return User{}, errors.New("User already exists")
	}
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
	err := collection.FindId(CreateIdFromString(guid)).One(&user)
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

func GetUsersCount() (int, error) {
	count, err := collection.Count()
	if err != nil {
		return 0, err
	}
	return count, err
}

func IsUserExisting(email string) (bool, error) {
	_, err := GetUserByEmail(email)
	if (err != nil) && (err.Error() == "not found") {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
