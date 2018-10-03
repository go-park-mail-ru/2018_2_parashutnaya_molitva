package user

import "testing"

type TestCaseValidationPassword struct {
	password string
	err      error
}

func TestValidationPassword(t *testing.T) {
	tests := []*TestCaseValidationPassword{
		&TestCaseValidationPassword{
			password: "",
			err:      errEmptyPassword,
		},
		&TestCaseValidationPassword{
			password: "asd",
			err:      errInvalidPassword,
		},
		&TestCaseValidationPassword{
			password: "привет",
			err:      errInvalidPassword,
		},
		&TestCaseValidationPassword{
			password: "Gasdf1aadРусский",
			err:      errInvalidPassword,
		},
		&TestCaseValidationPassword{
			password: "Gasdfgh1",
		},
	}

	for idx, test := range tests {
		err := ValidatePassword(test.password)

		if test.err == nil && err != nil {
			t.Fatalf("[%v] Unexpected error %#v", idx, err)
		}

		if test.err != err {
			t.Fatalf("[%v] %#v - was expected but got: %#v", idx, test.err, err)
		}
	}
}

type TestCaseValidationEmail struct {
	email string
	err   error
}

func TestValidationEmail(t *testing.T) {
	tests := []*TestCaseValidationEmail{
		&TestCaseValidationEmail{
			email: "",
			err:   errEmptyEmail,
		},
		&TestCaseValidationEmail{
			email: "asd",
			err:   errInvalidEmail,
		},
		&TestCaseValidationEmail{
			email: "@mail",
			err:   errInvalidEmail,
		},
		&TestCaseValidationEmail{
			email: "@",
			err:   errInvalidEmail,
		},
		&TestCaseValidationEmail{
			email: "@mail.ru",
			err:   errInvalidEmail,
		},
		&TestCaseValidationEmail{
			email: "asdsadasd@mail.ru",
		},
	}

	for idx, test := range tests {
		err := ValidateEmail(test.email)

		if test.err == nil && err != nil {
			t.Fatalf("[%v] Unexpected error %#v", idx, err)
		}

		if test.err != err {
			t.Fatalf("[%v] %#v - was expected but got: %#v", idx, test.err, err)
		}
	}
}
