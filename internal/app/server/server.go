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

type ServerData struct {
	port int
}

var (
	errNoPort = errors.New("Port wasn't passed")
)

func StartApp(port int) error {
	if port == -1 {
		return errNoPort
	}

	stringPort := ":" + strconv.Itoa(port)

	err := configReader.Read(confifFile, &corsData)
	if err != nil {
		singletoneLogger.LogError(err)
	}

	singletoneLogger.LogMessage("Server starting at " + stringPort)
	router := routes.NewRouter(http.DefaultServeMux)

	router.HandleFunc("/api/signin", middlewareChain(controllers.Foo, corsMiddleware))
	router.HandleFunc("/api/signup", middlewareChain(controllers.Foo, corsMiddleware))
	router.HandleFunc("/api/session/", middlewareChain(controllers.Session, corsMiddleware, authMiddleware)).Method("POST", "GET")
	router.HandleFunc("/api/user/:guid", middlewareChain(controllers.GetUser, corsMiddleware, authMiddleware)).Method("GET")

	// Документация
	router.HandleFunc("/docks/*", httpSwagger.WrapHandler)

	// Статика
	router.Handle("/storage/*", fileStorage.StorageHandler)

	return http.ListenAndServe(stringPort, router)
}
