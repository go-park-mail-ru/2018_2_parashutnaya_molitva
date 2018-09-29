package server

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
)

var (
	errNoPort = errors.New("Port wasn't passed")
)

func StartApp(port int) error {
	if port == -1 {
		return errNoPort
	}

	stringPort := ":" + strconv.Itoa(port)

	singletoneLogger.LogMessage("Server starting at " + stringPort)
	router := routes.NewRouter(http.DefaultServeMux)
	signIn := middlewareChain(controllers.SignIn, authMiddleware)
	router.HandleFunc("/api/session", signIn).Method("POST")
	router.HandleFunc("/api/signup", controllers.SignUp).Method("POST")
	router.HandleFunc("/user/:id/:kek", controllers.User)
	router.HandleFunc("/foo/*", controllers.Foo).Method("GET")

	return http.ListenAndServe(stringPort, router)
}
