package user

import (
	"github.com/globalsign/mgo"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/db"
)

var collection *mgo.Collection

func init() {
	collection = db.GetInstance().C("users")
}
