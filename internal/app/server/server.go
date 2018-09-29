package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
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
	sessionController := middlewareChain(controllers.Session, authMiddleware, corsMiddleware)
	router.HandleFunc("/api/session/", sessionController).Method("GET", "POST")

	router.HandleFunc("/user/", controllers.User)
	return http.ListenAndServe(stringPort, router)
}

type Test struct {
}

func (t *Test) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println("kek")
}
