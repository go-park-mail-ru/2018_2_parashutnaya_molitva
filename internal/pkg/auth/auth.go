package auth

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
)

func SetSession(guid string) (string, error) {
	token, err := updateOrAddIfNotExist(guid)
	if err != nil {
		singletoneLogger.LogError(err)
		return token, errors.New("Couldn't set the session")
	}
	return token, nil
}

func CheckSession(guid string, token string) (bool, error) {
	status, err := check(guid, token)
	if err != nil {
		singletoneLogger.LogError(err)

	}
	switch status {
	case statusNotExist:
		return false, errors.New("User not exist")
	case statusBadToken:
		return false, errors.New("Bad token")
	case statusExpired:
		return false, errors.New("Session expired")
	case statusOk:
		return true, nil
	default:
		return false, errors.New("Couldn't check the session")
	}
}


func DeleteSession(guid string) error {
	err := reset(guid)
	if err != nil {
		singletoneLogger.LogError(err)
		return errors.New("Couldn't delete session")
	}
	return nil
}
