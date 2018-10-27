package chess

type Bishop struct {
	Piece
}

func (p *Bishop) ShortName() rune {
	if p.color == WHITE {
		return 'B'
	} else {
		return 'b'
	}
}

func NewBishop(color PieceColor) *Bishop {
	return &Bishop{
		Piece{
			pieceType:BishopType,
			color:color,
			isMoved:false,
		},
	}
}