package chat

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	errNoRoom               = errors.New("No Room")
	errUnexpectedClose      = errors.New("Unexpected Close")
	errInitMsgWaitTooLong   = errors.New("Waiting Init Message too long")
	errInvalidMsgTypeInit   = errors.New("Invalid type: wating for Init")
	errInvalidMsgInitFormat = errors.New("Invalid init message format")
	errInavlidMsgFormat     = errors.New("Invalid message format")
	errAlreadySearching     = errors.New("Already searching")
	errAlreadyInRoom        = errors.New("Already in room")
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
