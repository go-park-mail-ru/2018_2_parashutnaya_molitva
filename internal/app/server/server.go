package server

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

var (
	errNoPort = errors.New("Port wasn't passed")
)

func StartApp(port int) error {
	if port == -1 {
		return errNoPort
	}

	stringPort := ":" + strconv.Itoa(port)

	singletoneLogger.LogMessage(stringPort)
	router := routes.NewRouter(http.DefaultServeMux)
	router.HandleFunc("/api/signin", controllers.SignIn).Method("GET")
	router.HandleFunc("/api/signup", controllers.SignUp).Method("POST")
	return http.ListenAndServe(stringPort, router)
}