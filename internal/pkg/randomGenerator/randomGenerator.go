package randomGenerator

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

func RandomBytes(count int) ([]byte, error) {
	bytes := make([]byte, count)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func RandomString(length int) (string, error) {
	bytes, err := RandomBytes(length)
	return base64.URLEncoding.EncodeToString(bytes), err
}

func RandomBool() bool {
	c := make(chan struct{})
	close(c)
	select {
	case <-c:
		return true
	case <-c:
		return false
	}
}

func RandomUUID() (string, error) {
	id, err := uuid.NewRandom()
	return id.String(), err
}
