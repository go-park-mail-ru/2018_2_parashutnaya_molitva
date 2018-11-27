package user

import (
	simpleErrors "errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
)

//easyjson:json
type User struct {
	Guid         bson.ObjectId `bson:"_id" json:"guid"`
	Login        string        `bson:"login" json:"login"`
	Email        string        `bson:"email" json:"email"`
	Password     string        `bson:"-" json:"-"`
	HashPassword string        `bson:"password" json:"-"`
	Avatar       string        `bson:"avatar" json:"avatar"`
	Score        int           `bson:"score" json:"score"`
}

func (u *User) ChangeAvatar(avatarName string) error {
	u.Avatar = avatarName
	err := collection.UpdateId(u.Guid, u)
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	return nil
}

func SigninUser(loginOrEmail string, password string) (User, error) {
	var err error
	var u User
	if emailRegex.MatchString(loginOrEmail) {
		u, err = GetUserByEmail(loginOrEmail)
	} else {
		u, err = GetUserByLogin(loginOrEmail)
	}
	if (err != nil) && (err.Error() == "not found") {
		return User{}, err//simpleErrors.New("userController not found")
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

func (u *User) ChangeLogin(login string) error {
	u.Login = login
	err := collection.UpdateId(u.Guid, u)
	return err
}

func (u *User) ChangeEmail(email string) error {
	u.Email = email
	err := collection.UpdateId(u.Guid, u)
	return err
}

func (u *User) ChangePassword(password string) error {
	hashedPassword, err := hashPassword(password)

	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	u.HashPassword = hashedPassword
	err = collection.UpdateId(u.Guid, u)
	return err
}

func (u *User) AddScore(score int) error {
	u.Score += score
	err := collection.UpdateId(u.Guid, u)
	return errors.WithStack(err)
}

func CreateUser(email string, login string, password string) (User, error) {
	isEmailExisting, err := IsUserEmailExisting(email)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, simpleErrors.New("Unknown error")
	}
	if isEmailExisting {
		return User{}, simpleErrors.New("userController with such email already exists")
	}
	isLoginExisting, err := IsUserLoginExisting(login)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, simpleErrors.New("Unknown error")
	}
	if isLoginExisting {
		return User{}, simpleErrors.New("userController with such login already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		singletoneLogger.LogError(err)
		return User{}, err
	}
	u := User{
		Guid:         bson.NewObjectId(),
		Login:        login,
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

func (u *User) UpdateUser(updateUser UpdateUserStruct) error {
	if updateUser.Avatar != nil {
		u.Avatar = updateUser.Avatar.(string)
	}
	//if updateUser.Email != nil {
	//	u.Email = updateUser.Email.(string)
	//}
	if updateUser.Password != nil {
		hashedPassword, err := hashPassword(updateUser.Password.(string))
		if err != nil {
			return err
		}
		u.HashPassword = hashedPassword
	}
	err := collection.UpdateId(u.Guid, u)
	return err
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
	err := collection.Find(bson.M{"email": email}).Collation(getEmailCollation()).One(&user)
	return user, err
}

func GetUserByLogin(login string) (User, error) {
	user := User{}
	err := collection.Find(bson.M{"login": login}).One(&user)
	return user, err
}

func CreateIdFromString(str string) bson.ObjectId {
	if !bson.IsObjectIdHex(str) {
		return ""
	}
	objId := bson.ObjectIdHex(str)
	return objId
}

func GetUsersCount() (int, error) {
	count, err := collection.Count()
	if err != nil {
		return 0, err
	}
	return count, err
}

func IsUserEmailExisting(email string) (bool, error) {
	_, err := GetUserByEmail(email)
	if (err != nil) && (err.Error() == "not found") {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func IsUserLoginExisting(login string) (bool, error) {
	_, err := GetUserByLogin(login)
	if (err != nil) && (err.Error() == "not found") {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

var (
	errTypeCastEmail    = errors.New("Can't cast Email interface{} to string")
	errTypeCastPassword = errors.New("Can't cast Password interface{} to string")
)

//easyjson:json
type UpdateUserStruct struct {
	Avatar interface{} `json:"avatar"`
	Login  interface{} `json:"login"`
	//Email    interface{} `json:"email"`
	Password interface{} `json:"password"`
}

func (u *UpdateUserStruct) Validate() (string, error) {

	//if u.Email != nil {
	//	email, ok := u.Email.(string)
	//	if !ok {
	//		singletoneLogger.LogError(errTypeCastEmail)
	//		return "", errTypeCastEmail
	//	}
	//
	//	if err := ValidateEmail(email); err != nil {
	//		return "email", err
	//	}
	//}

	if u.Password != nil {
		password, ok := u.Password.(string)
		if !ok {
			singletoneLogger.LogError(errTypeCastPassword)
			return "", errTypeCastPassword
		}

		if err := ValidatePassword(password); err != nil {
			return "password", err
		}
	}

	if u.Login != nil {
		login, ok := u.Login.(string)
		if !ok {
			err := errors.New("Can't cast Login interface{} to string")
			singletoneLogger.LogError(err)
			return "", err
		}

		if err := ValidateLogin(login); err != nil {
			return "password", err
		}
	}

	return "", nil

}

func getEmailCollation() *mgo.Collation {
	return &mgo.Collation{Locale: "en", Strength: 1}
}

func getLoginCollation() *mgo.Collation {
	return getEmailCollation()
}
