package gamelogic

type ChessEngine interface {
	Move(string) error
	IsCheckmate() bool
	IsStalemate() bool
}
