package chess

type Pawn struct {
	Piece
}

func (p *Pawn) ShortName() rune {
	if p.color == WHITE {
		return 'P'
	} else {
		return 'p'
	}
}

func NewPawn(color PieceColor) *Pawn {
	return &Pawn{
		Piece{
			pieceType: PawnType,
			color:     color,
			isMoved:   false,
		},
	}
}
