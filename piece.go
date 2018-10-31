package chess

type Piece struct {
	pieceType PieceType
	color     PieceColor
	isMoved   bool
	PieceInterface
}

type PieceInterface interface {
	Type() PieceType
	Color() PieceColor
	IsMoved() bool
	SetMoved(bool)
	ShortName() rune
}

func (p *Piece) Type() PieceType {
	return p.pieceType
}

func (p *Piece) Color() PieceColor {
	return p.color
}

func (p *Piece) IsMoved() bool {
	return p.isMoved
}

func (p *Piece) SetMoved(bool) {
	p.isMoved = true
}
