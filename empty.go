package chess

type Empty struct {
	Piece
}

func (p *Empty) ShortName() rune {
	return '.'
}

func NewEmpty() *Empty {
	return &Empty{
		Piece{
			pieceType: EmptyType,
			color:     NONE,
			isMoved:   false,
		},
	}
}
