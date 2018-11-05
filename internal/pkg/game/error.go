package game

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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

func sendCloseError(conn *websocket.Conn, code int, errorRaw string) {
	msg := websocket.FormatCloseMessage(code, errorRaw)
	err := conn.WriteMessage(websocket.CloseMessage, msg)
	if err != nil {
		singletoneLogger.LogError(errors.WithStack(err))
	}
}
