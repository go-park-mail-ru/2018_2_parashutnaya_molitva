package mainServer

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/auth"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)
type Server struct {

}

func GRPCServer () {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}
	grpcServer := grpc.NewServer()
	RegisterAuthCheckerServer(grpcServer, &Server{})

	singletoneLogger.LogMessage("start serve authChecker in 8081 port")
	grpcServer.Serve(lis)
}

func (s *Server) AuthUser(ctx context.Context, session *Session) (*User, error) {
	cookie := session.Cookie
	isAuth, guid, err := auth.CheckSession(cookie)
	if err != nil {
		singletoneLogger.LogError(err)
		return nil, err
	}
	uData, err := user.GetUserByGuid(guid)
	if err != nil {
		singletoneLogger.LogError(err)
		return nil, err
	}
	u := &User{
		IsAuth: isAuth,
		Guid:guid,
		Login:uData.Login,
		Email:uData.Email,
	}
	return u, nil
}