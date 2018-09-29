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
	router.HandleFuncWithMiddleware("/api/session/", controllers.Session).Method("POST", "GET")
	router.HandleFuncWithMiddleware("/api/user/:id", controllers.GetUser).Method("GET")
	router.HandleFuncWithMiddleware("/api/user/:id", controllers.ChangeUser).Method("PUT")
	router.HandleFuncWithMiddleware("/api/user/:id", controllers.DeleteUser).Method("DELETE")
	router.HandleFuncWithMiddleware("/api/user/", controllers.SaveUser).Method("POST")

	return http.ListenAndServe(stringPort, router)
}
