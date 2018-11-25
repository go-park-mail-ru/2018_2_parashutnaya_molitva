package user

//easyjson:json
type UserScore struct {
	Login string `json:"login" bson:'login'`
	Score int    `json:"score" bson:"score"`
}

func GetScores(limit int, offset int) ([]UserScore, error) {
	var result []UserScore
	err := collection.Find(nil).Sort("-score").Skip(offset).Limit(limit).All(&result)
	return result, err
}
