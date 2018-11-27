package game

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/game"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC/core"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func StartGameServer(port int) error {
	portString := strconv.Itoa(port)
	grcpConn, err := grpc.Dial(
		"127.0.0.1:" + gRPC.GetCorePort(),
		grpc.WithInsecure(),
	)
	if err != nil {
		singletoneLogger.LogError(errors.New("cant connect to grpc"))
	}
	defer grcpConn.Close()

	gameServer := game.NewGame(core.NewCoreClient(grcpConn))

	router := mux.NewRouter()
	router.Use()
	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.RecoverPanicMiddleware)
	router.Handle("/api/game", &controllers.FindRoom{Game: gameServer}).Methods("POST", "OPTIONS")
	router.Handle("/api/game/ws", &controllers.StartGame{Game: gameServer, Upgrader: &websocket.Upgrader{}})
	fmt.Printf("starting chat server at %s\n", portString)
	return http.ListenAndServe(":"+ portString, router)
}
