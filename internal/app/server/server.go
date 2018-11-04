package server

import (
	"net/http"
	"strconv"

	_ "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/docs"
	"github.com/gorilla/websocket"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/fileStorage"
	g "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/game"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger"
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

	router := mux.NewRouter()
	router.Use(authMiddleware)
	router.Use(corsMiddleware)

	router.HandleFunc("/api/session", controllers.Session).Methods("POST", "GET", "OPTIONS")
	router.HandleFunc("/api/session", controllers.DeleteSession).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/user/count", controllers.GetUsersCount).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/score/", controllers.GetUsersScore).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{guid}", controllers.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{guid}", controllers.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/avatar", controllers.UploadAvatar).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user", controllers.CreateUser).Methods("POST", "OPTIONS")

	game := g.NewGame()
	router.Handle("/api/game/", &controllers.FindRoom{game}).Methods("POST", "OPTIONS")

	router.Handle("/api/game/ws/", &controllers.StartGame{game, &websocket.Upgrader{}})
	// Документация
	router.HandleFunc("/docs/*", httpSwagger.WrapHandler)
	router.PathPrefix("/storage/").Handler(fileStorage.StorageHandler)

	return http.ListenAndServe(stringPort, router)
}
