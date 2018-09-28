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
	signin := middlewareChain(controllers.SignIn, auth)
	router.HandleFunc("/api/signin", signin).Method("GET")
	router.HandleFunc("/api/signup", controllers.SignUp).Method("POST")
	router.HandleFunc("/user/:id", controllers.User)
	router.HandleFunc("/user/", controllers.User)
	return http.ListenAndServe(stringPort, router)
}
