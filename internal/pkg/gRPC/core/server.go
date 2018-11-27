package core

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/auth"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/user"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type Server struct {

}

func GRPCServer() {
	port := gRPC.GetCorePort()
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		singletoneLogger.LogError(err)
		return
	}
	grpcServer := grpc.NewServer()
	RegisterCoreServer(grpcServer, &Server{})

	singletoneLogger.LogMessage(fmt.Sprintf("start serve core grpc in %s port", port))
	grpcServer.Serve(lis)
}

func (*Server) GetUserBySession(ctx context.Context, session *Session) (*User, error) {
	cookie := session.Cookie
	isAuth, guid, err := auth.CheckSession(cookie)
	if err != nil {
		singletoneLogger.LogError(err)
		return nil, err
	}
	if !isAuth {
		return &User{
			IsAuth: false,
		}, nil
	}
	uData, err := user.GetUserByGuid(guid)
	if err != nil {
		singletoneLogger.LogError(err)
		return nil, err
	}
	return &User{
		IsAuth: isAuth,
		Guid:   guid,
		Login:  uData.Login,
		Email:  uData.Email,
		Score:  int32(uData.Score),
	}, nil
}

func (*Server) AddScoreToUser(ctx context.Context, scoreAdd *ScoreAdd) (*Nothing, error) {
	u, err := user.GetUserByGuid(scoreAdd.Guid)
	if err != nil {
		return nil, errors.Wrap(err, "can't find user")
	}
	err = u.AddScore(int(scoreAdd.Score))
	if err != nil {
		return nil, errors.Wrap(err, "can't update score of user")
	}
	return nil, nil
}
