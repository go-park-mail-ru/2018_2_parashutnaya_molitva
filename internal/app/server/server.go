package server

import (
	"net/http"
	"strconv"

	_ "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/fileStorage"
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
	router.HandleFuncWithMiddleware("/api/user/:guid", controllers.GetUser).Method("GET")
	//router.HandleFuncWithMiddleware("/api/user/:guid", controllers.ChangeUser).Method("PUT")
	//router.HandleFuncWithMiddleware("/api/user/:guid", controllers.DeleteUser).Method("DELETE")
	//router.HandleFuncWithMiddleware("/api/user/", controllers.SaveUser).Method("POST")

	// Документация
	router.HandleFunc("/docks/*", httpSwagger.WrapHandler)

	// Статика
	router.HandleFunc("/storage/*", fileStorage.StorageHandleFunc)

	return http.ListenAndServe(stringPort, router)
}
