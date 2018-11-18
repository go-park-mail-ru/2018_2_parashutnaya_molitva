package auth

import (
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
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
	TestUserToStay.Token, _, err = SetSession(TestUserToStay.Guid)
	if err != nil {
		t.Error(err)
	}
	TestUserToDelete.Token, _, err = SetSession(TestUserToDelete.Guid)
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
	status, guid, err := CheckSession(TestUserToDelete.Token)
	if (status != false) || (err == nil) || (guid == TestUserToDelete.Guid) {
		t.Error(errors.Wrap(err, "Sees deleted session as Valid"))
	}
}

func TestCheckValidSession(t *testing.T) {
	status, guid, err := CheckSession(TestUserToStay.Token)
	if (status != true) || (err != nil) || (guid != TestUserToStay.Guid) {
		t.Error(errors.Wrap(err, "Sees valid session as invalid"))
	}
}

func TestCheckSessionWithWrongToken(t *testing.T) {
	randomToken, err := generateToken()
	if err != nil {
		t.Error(err)
	}
	status, guid, err := CheckSession(randomToken)
	if (status != false) || (err == nil) || (guid == TestUserToStay.Guid) {
		t.Error(errors.Wrap(err, "Sees invalid session as valid"))
	}
}
