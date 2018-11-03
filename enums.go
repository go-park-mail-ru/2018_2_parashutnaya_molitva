package chess

type pieceColor int

const (
	white pieceColor = 1
	black pieceColor = 0
	none  pieceColor = -1
)

type pieceType int

const (
	emptyType pieceType = iota
	noneType
	enPassantType
	pawnType
	knightType
	bishopType
	rookType
	queenType
	kingType
)
