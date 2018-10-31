package chess

type King struct {
	Piece
}

func (p *King) ShortName() rune {
	if p.color == WHITE {
		return 'K'
	} else {
		return 'k'
	}
}

func NewKing(color PieceColor) *King {
	return &King{
		Piece{
			pieceType: KingType,
			color:     color,
			isMoved:   false,
		},
	}
}
