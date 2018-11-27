package auth

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/auth"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type Auth struct {

}

func GRPCServer() error{
	port := gRPC.GetAuthPort()
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		singletoneLogger.LogError(err)
		return err
	}
	grpcServer := grpc.NewServer()
	RegisterAuthServer(grpcServer, &Auth{})

	singletoneLogger.LogMessage(fmt.Sprintf("start serve auth grpc in %s port", port))
	return grpcServer.Serve(lis)
}

func (*Auth) DeleteSession(ctx context.Context, session  *SessionCookie) (*Nothing, error) {
	err := auth.DeleteSession(session.Cookie)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (*Auth) CheckSession(ctx context.Context, session *SessionCookie) (*UserAuthed, error) {
	isAuth, guid, err := auth.CheckSession(session.Cookie)
	if err != nil {
		return nil, err
	}
	return &UserAuthed{
		IsAuth: isAuth, Guid: guid,
	}, nil
}

func (*Auth) SetSession(ctx context.Context, userGuid *UserGuid) (*SessionToken, error) {
	token, timeExpire, err := auth.SetSession(userGuid.Guid)
	if err != nil {
		return nil, err
	}
	return &SessionToken{
		Token:token,
		TimeExpire:timeExpire.String(),
	}, nil
}