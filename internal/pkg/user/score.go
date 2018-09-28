package user

type UserScore struct {
	Email string `bson:"email"`
	Score int    `bson:"score"`
}

type GetScoreParameters struct {
	Limit  int
	Offset int
}

func GetScores(param GetScoreParameters) ([]UserScore, error) {
	var result []UserScore
	err := collection.Find(nil).Skip(param.Offset).Limit(param.Limit).All(&result)
	return result, err
}
