package chess

type Queen struct {
	Piece
}

func (p *Queen) ShortName() rune {
	if p.color == WHITE {
		return 'Q'
	} else {
		return 'q'
	}
}

func NewQueen(color PieceColor) *Queen {
	return &Queen{
		Piece{
			pieceType: QueenType,
			color:     color,
			isMoved:   false,
		},
	}
}
