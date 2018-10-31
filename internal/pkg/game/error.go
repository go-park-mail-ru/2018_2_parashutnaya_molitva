package game

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
)

func sendError(conn *websocket.Conn, errorRaw string) {

	msg, err := MarshalToMessage(ErrorMsg, &ErrorMessage{errorRaw})
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}

	err = conn.WriteJSON(msg)
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}
}
