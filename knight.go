package chess

type Knight struct {
	Piece
}

func (n *Knight) ShortName() rune {
	if n.color == WHITE {
		return 'N'
	} else {
		return 'n'
	}
}

func NewKnight(color PieceColor) *Knight {
	return &Knight{
		Piece{
			pieceType:KnightType,
			color:color,
			isMoved:false,
		},
	}
}