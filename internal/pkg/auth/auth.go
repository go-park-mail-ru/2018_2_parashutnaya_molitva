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

func CheckSession(token string) (bool, string, error) {
	status, guid, err := check(token)
	if err != nil {
		singletoneLogger.LogError(err)

	}
	switch status {
	case statusBadToken:
		return false, "", errors.New("Bad token")
	case statusExpired:
		return false, "", errors.New("Session expired")
	case statusOk:
		return true, guid, nil
	default:
		return false, "", errors.New("Couldn't check the session")
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
