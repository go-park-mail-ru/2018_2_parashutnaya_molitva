package user

import (
	"regexp"
	"unicode"

	"github.com/pkg/errors"
)

var (
	errEmptyPassword   = errors.New("Empty password")
	errEmptyEmail      = errors.New("Empty email")
	errEmptyLogin   = errors.New("Empty login")
	errLoginTooLong   = errors.New("Too long login")
	errInvalidEmail    = errors.New("Email is invalid")
	errInvalidPassword = errors.New("Must contain at least 8 characters, 1 number, 1 upper and 1 lowercase")
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
	var number, upper, lower bool
	ascii := true

	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			ascii = false
			return false
		}

		switch {
		case unicode.IsNumber(r):
			number = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		}
	}

	return (len(s) >= 8) && number && upper && lower && ascii
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
