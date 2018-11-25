package chat

import (
	"fmt"
	"net/http"

	chatModel "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/chat"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC/mainServer"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func StartChatServer() error {
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		singletoneLogger.LogError(errors.New("cant connect to grpc"))
	}
	defer grcpConn.Close()

	chat := chatModel.NewChat(mainServer.NewAuthCheckerClient(grcpConn))

	router := mux.NewRouter()
	router.Handle("/api/chat/ws", &chatModel.StartChat{Chat: chat, Upgrader: &websocket.Upgrader{}})
	fmt.Println("starting chat server at :3335")
	return http.ListenAndServe(":3335", router)
}
