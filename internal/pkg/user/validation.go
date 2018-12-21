package user

import (
	"github.com/pkg/errors"
	"regexp"
)

var (
	errEmptyPassword   = errors.New("Empty password")
	errEmptyEmail      = errors.New("Empty email")
	errEmptyLogin   = errors.New("Empty login")
	errLoginTooLong   = errors.New("Too long login")
	errInvalidEmail    = errors.New("Email is invalid")
	errInvalidPassword = errors.New("Must contain at least 4 characters")
)

var (
	// RFC 2822
	emailRegex = regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
)

func ValidatePassword(password string) error {
	if password == "" {
		return errEmptyPassword
	}

	if !verifyPassword(password) {
		return errInvalidPassword
	}

	return nil
}

func ValidateLogin(login string) error {
	if login == "" {
		return errEmptyLogin
	}
	if len(login) > 12 {
		return errLoginTooLong
	}
	return nil
}

func verifyPassword(s string) bool {
	return len(s) >= 4
}

func ValidateEmail(email string) error {
	if email == "" {
		return errEmptyEmail
	}

	if !emailRegex.MatchString(email) {
		return errInvalidEmail
	}

	return nil
}
