package auth

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type TestUserStruct struct {
	Guid  string
	Token string
}

var TestUserToStay = TestUserStruct{
	Guid: bson.NewObjectId().Hex(),
}

var TestUserToDelete = TestUserStruct{
	Guid: bson.NewObjectId().Hex(),
}

func TestSetSession(t *testing.T) {
	var err error
	TestUserToStay.Token, err = SetSession(TestUserToStay.Guid)
	if err != nil {
		t.Error(err)
	}
	TestUserToDelete.Token, err = SetSession(TestUserToDelete.Guid)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteSession(t *testing.T) {
	err := DeleteSession(TestUserToDelete.Guid)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckDeletedSession(t *testing.T) {
	status, err := CheckSession(TestUserToDelete.Guid, TestUserToDelete.Token)
	if (status != false) || (err == nil) {
		t.Error(errors.Wrap(err, "Sees deleted session as Valid"))
	}
}

func TestCheckValidSession(t *testing.T) {
	status, err := CheckSession(TestUserToStay.Guid, TestUserToStay.Token)
	if (status != true) || (err != nil) {
		t.Error(errors.Wrap(err, "Sees valid session as invalid"))
	}
}

func TestCheckSessionWithWrongToken(t *testing.T) {
	randomToken, err := generateToken()
	if err != nil {
		t.Error(err)
	}
	status, err := CheckSession(TestUserToStay.Guid, randomToken)
	if status != false || err == nil {
		t.Error(errors.Wrap(err, "Sees invalid session as valid"))
	}
}
