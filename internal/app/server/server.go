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
	router.Use(authMiddleware)
	router.Use(corsMiddleware)
	router.HandleFunc("/api/session/", controllers.Session).Method("POST", "GET", "OPTIONS")
	router.HandleFunc("/api/session/", controllers.DeleteSession).Method("DELETE", "OPTIONS")
	router.HandleFunc("/api/user/count/", controllers.GetUsersCount).Method("GET", "OPTIONS")
	router.HandleFunc("/api/user/score/*", controllers.GetUsersScore).Method("GET", "OPTIONS")
	router.HandleFunc("/api/user/:guid", controllers.GetUser).Method("GET", "OPTIONS")
	router.HandleFunc("/api/user/:guid", controllers.UpdateUser).Method("PUT", "OPTIONS")
	router.HandleFunc("/api/avatar/", controllers.UploadAvatar).Method("POST", "OPTIONS")
	router.HandleFunc("/api/user/", controllers.CreateUser).Method("POST", "OPTIONS")

	// Документация
	router.HandleFunc("/docks/*", httpSwagger.WrapHandler)

	// Статика
	router.Handle("/storage/*", fileStorage.StorageHandler)

	return http.ListenAndServe(stringPort, router)
}
