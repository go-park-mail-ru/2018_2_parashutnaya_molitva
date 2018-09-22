package auth

type userAuthData struct {
	guid     int // global user id
	login    string
	password string
}

//func (user *userAuthData) registerUser() userAuthData {
//	return
//}
//
//func (user *userAuthData) authUser() (sessionCookie string) {
//	return ""
//}
