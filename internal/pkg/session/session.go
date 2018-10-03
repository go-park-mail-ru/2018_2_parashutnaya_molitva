package session

import (
	"net/http"
	"time"
)

const CookieName = "sessionid"

func CreateAuthCookie(token string, timeExpire time.Time) *http.Cookie {
	return &http.Cookie{
		Name:    CookieName,
		Value:   token,
		Expires: timeExpire,
	}
}
