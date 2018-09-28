package randomGenerator

import (
	"crypto/rand"
	"encoding/base64"
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
