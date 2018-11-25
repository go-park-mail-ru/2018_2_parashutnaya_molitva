package chat

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC/mainServer"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/session"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"net/http"
)

type Chat struct {
	authService mainServer.AuthCheckerClient
}

func NewChat(authService mainServer.AuthCheckerClient) *Chat {
	return &Chat{
		authService,
	}
}

type StartChat struct {
	Chat     *Chat
	Upgrader *websocket.Upgrader
}

func (sc *StartChat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userData := &mainServer.User{}
	cookie, noCookie := r.Cookie(session.CookieName)
	if noCookie != nil {
		singletoneLogger.LogMessage(noCookie.Error())
	} else {
		var err error
		userData, err = sc.Chat.authService.AuthUser(ctx, &mainServer.Session{Cookie:cookie.Value})
		if err != nil {
			singletoneLogger.LogError(err)
		}
	}

	// Уже есть CORS middleware, который отклоняет запросы с неразрешенных Origin
	sc.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := sc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}

	go sc.Chat.initConnection(userData, conn)
}

func (c *Chat) initConnection (user *mainServer.User, conn *websocket.Conn) {
// ВЛАД ПИШИ ТУТ
}