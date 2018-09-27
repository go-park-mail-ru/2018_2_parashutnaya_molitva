package user

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var collection *mgo.Collection

func init() {
	collection = db.GetInstance().C("users")
}

func GetUserByGuid(guid string) (*User, error) {
	user := &User{}
	err := collection.FindId(bson.M{"_id": bson.ObjectIdHex(guid)}).One(&user)
	return user, err
}

func (u *User) updateUser() error {
	err := collection.UpdateId(bson.M{"_id": u.Guid}, u)
	return err
}

func (u *User) deleteUser() error {
	err := collection.RemoveId(bson.M{"_id": u.Guid})
	return err
}

func deleteUserById(guid string) error {
	err := collection.RemoveId(bson.M{"_id": guid})
	return err
}

func getUserByEmail(email string) (*User, error) {
	user := &User{}
	err := collection.Find(bson.M{"email": email}).One(&user)
	return user, err
}
