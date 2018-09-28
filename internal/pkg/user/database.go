package user

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/db"
	"gopkg.in/mgo.v2"
)

var collection *mgo.Collection

func init() {
	collection = db.GetInstance().C("users")
}
