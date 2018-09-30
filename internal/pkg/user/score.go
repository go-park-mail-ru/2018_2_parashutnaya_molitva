package user

//easyjson:json
type UserScore struct {
	Email string `json:"email" bson:"email"`
	Score int    `json:"score" bson:"score"`
}

func GetScores(	limit  int, offset int) ([]UserScore, error) {
	var result []UserScore
	err := collection.Find(nil).Sort("-score").Skip(offset).Limit(limit).All(&result)
	return result, err
}
