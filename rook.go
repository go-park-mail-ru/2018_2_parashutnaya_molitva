package chess

type Rook struct {
	Piece
}

func (p *Rook) ShortName() rune {
	if p.color == WHITE {
		return 'R'
	} else {
		return 'r'
	}
}

func NewRook(color PieceColor) *Rook {
	return &Rook{
		Piece{
			pieceType: RookType,
			color:     color,
			isMoved:   false,
		},
	}
}
