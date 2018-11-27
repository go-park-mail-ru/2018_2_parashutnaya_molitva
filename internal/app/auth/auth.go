package auth

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/gRPC/auth"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"time"
)

func StartAuth() {
	err := auth.GRPCServer()
	if err != nil {
		singletoneLogger.LogError(err)
		time.Sleep(1)
	}
}
