package chess

type None struct {
	Piece
}

func NewNone() *None {
	n := &None{}
	n.pieceType = NoneType

	return n
}

func (p *None) ShortName() rune {
	return '!'
}
