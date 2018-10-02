package user

import "github.com/pkg/errors"

var (
	errEmptyPassword = errors.New("Empty password")
	errEmptyEmail    = errors.New("Empty email")
)

func ValidatePassword(password string) error {
	if password == "" {
		return errEmptyPassword
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errEmptyEmail
	}

	return nil
}
