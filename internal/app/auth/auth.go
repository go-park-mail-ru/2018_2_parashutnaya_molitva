package auth

import "errors"

func StartAuth(errChan chan<- error) {
	errChan <- errors.New("Auth wasn't written. Nikita Forever")
}
