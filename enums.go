package chess

type PieceColor int
const (
	WHITE PieceColor = 1
	BLACK PieceColor = 0
	NONE  PieceColor = -1
)

type PieceType int
const (
	EmptyType PieceType = iota
	NoneType
	EnPassantType
	PawnType
	KnightType
	BishopType
	RookType
	QueenType
	KingType
)
