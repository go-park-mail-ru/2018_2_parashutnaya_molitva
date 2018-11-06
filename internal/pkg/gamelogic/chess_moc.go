package gamelogic

import "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"

type ChessEngineMock struct{}

func (c *ChessEngineMock) Move(uci string) error {
	singletoneLogger.LogMessage("<ChessEngineMock>Move: " + uci)

	return nil
}

func (c *ChessEngineMock) IsCheckmate() bool {
	singletoneLogger.LogMessage("<ChessEngineMock>Check chessmate")
	return false
}

func (c *ChessEngineMock) IsStalemate() bool {
	singletoneLogger.LogMessage("<ChessEngineMock>Check stalemate")
	return false
}
