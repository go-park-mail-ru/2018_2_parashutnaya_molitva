package gamelogic

type ChessEngine interface {
	Move(string) error
	IsCheckmate() bool
	IsStalemate() bool
	IsInsufficientMaterial() bool
	IsGameOver() bool
	PrintLegalMoves()
	PrintBoard()
}
