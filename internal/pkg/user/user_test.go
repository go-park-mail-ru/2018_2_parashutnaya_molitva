package user

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"
)

type TestUserDataStruct struct {
	Email    string
	Login string
	Password string
}

var TestUserData1 = TestUserDataStruct{
	"test@mail.ru",
	"tes228",
	"1234",
}
var testUser *User

func TestIsUserExisting(t *testing.T) {
	isExisting, err := IsUserEmailExisting(TestUserData1.Email)
	if err != nil {
		t.Fatal(err)
	}
	if isExisting {
		t.Fatal(errors.New("userController should not exist"))
	}
}

func TestCreateUser(t *testing.T) {
	u, err := CreateUser(TestUserData1.Email, TestUserData1.Login, TestUserData1.Password)
	if err != nil {
		t.Fatal(err)
	}
	testUser = &u
}

func TestFindUserByEmail(t *testing.T) {
	u, err := GetUserByEmail(TestUserData1.Email)
	if err != nil {
		t.Fatal(err)
	}
	testUser = &u
}

func TestFindUserByLogin(t *testing.T) {
	u, err := GetUserByLogin(TestUserData1.Login)
	if err != nil {
		t.Fatal(err)
	}
	testUser = &u
}

func TestFindUserByEmailUppercase(t *testing.T) {
	u, err := GetUserByEmail(strings.ToUpper(TestUserData1.Email))
	if err != nil {
		t.Fatal(err)
	}
	testUser = &u
}

func TestLoginUser(t *testing.T) {
	u, err := SigninUser(TestUserData1.Email, TestUserData1.Password)
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
