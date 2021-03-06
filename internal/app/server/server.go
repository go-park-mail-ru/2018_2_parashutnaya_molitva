package server

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC/core"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/middleware"
	"net/http"
	"strconv"

	_ "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/docs"
	"github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/fileStorage"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type ServerData struct {
	port int
}

func StartApp(port int) error {
	if port == -1 {
		return errors.New("Port wasn't passed")
	}

	stringPort := ":" + strconv.Itoa(port)

	singletoneLogger.LogMessage("Server starting at " + stringPort)

	router := mux.NewRouter()
	router.Use(authMiddleware)
	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.RecoverPanicMiddleware)

	router.HandleFunc("/api/session", controllers.Session).Methods("POST", "GET", "OPTIONS")
	router.HandleFunc("/api/session", controllers.DeleteSession).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/user/count", controllers.GetUsersCount).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/score/", controllers.GetUsersScore).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{guid}", controllers.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{guid}", controllers.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/avatar", controllers.UploadAvatar).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user", controllers.CreateUser).Methods("POST", "OPTIONS")

	// Документация
	router.HandleFunc("/docs/*", httpSwagger.WrapHandler)
	router.PathPrefix("/storage/").Handler(fileStorage.StorageHandler)

	go core.GRPCServer()
	return http.ListenAndServe(stringPort, router)
}