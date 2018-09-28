package user

import "gopkg.in/mgo.v2/bson"

type UserScore struct {
	Guid         bson.ObjectId `bson:"_id"`
	Email string `bson:"email"`
	Score uint `bson:"score"`
}

type GetScoreParameters struct  {

}

func GetScore() {

}
