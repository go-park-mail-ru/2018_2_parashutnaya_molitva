package auth

type userAuthData struct {
	guid int // global user id
	login string
	password string
}

func (user* userAuthData) registerUser () userAuthData {

}

func (user* userAuthData) authUser () (sessionCookie string)  {
	sessionCookie = ""
	return
}