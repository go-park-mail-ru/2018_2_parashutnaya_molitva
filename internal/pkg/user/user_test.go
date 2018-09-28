package user

import (
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

type TestUserDataStruct struct {
	Email    string
	Password string
}

var TestUserData1 = TestUserDataStruct{
	"testтест@test.test",
	"1234",
}
var testUser *User

func TestIsUserExisting(t *testing.T) {
	isExisting, err := IsUserExisting(TestUserData1.Email)
	if err != nil {
		t.Fatal(err)
	}
	if isExisting {
		t.Fatal(errors.New("User should not exist"))
	}
}

func TestCreateUser(t *testing.T) {
	u, err := CreateUser(TestUserData1.Email, TestUserData1.Password)
	if err != nil {
		t.Fatal(err)
	}
	testUser = &u
}

func TestFindUser(t *testing.T) {
	u, err := GetUserByEmail(TestUserData1.Email)
	if err != nil {
		t.Fatal(err)
	}
	testUser = &u
}

func TestLoginUser(t *testing.T) {
	u, err := LoginUser(TestUserData1.Email, TestUserData1.Password)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(u, *testUser) {
		t.Fatal(errors.New("is not the same user"))
	}
}

func TestDeleteUser(t *testing.T) {
	err := testUser.DeleteUser()
	if err != nil {
		t.Fatal(errors.New("cant delete user"))
	}
}
