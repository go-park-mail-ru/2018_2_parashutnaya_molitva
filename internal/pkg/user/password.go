package user

import (
	"golang.org/x/crypto/bcrypt"
)

var encriptionCost = 15

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), encriptionCost)
	return string(bytes), err
}

func checkPasswordByHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
