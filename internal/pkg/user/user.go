package user

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Guid     bson.ObjectId `bson:"_id"` // global user id
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
	Avatar   string
}
