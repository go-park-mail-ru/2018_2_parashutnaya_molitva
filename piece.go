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

func NewPiece(pt PieceType, color PieceColor) Piece {
	return Piece{
		pieceType: pt,
		color:     color,
		isMoved:   false,
	}
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

func (p *Piece) ShortName() rune {

	switch p.Type() {
	case EmptyType:
		{
			return '.'
		}
	case NoneType:
		{
			return '!'
		}
	case EnPassantType:
		{
			if p.Color() == WHITE {
				return 'E'
			}
			return 'e'
		}
	case PawnType:
		{
			if p.Color() == WHITE {
				return 'P'
			}
			return 'p'
		}
	case KnightType:
		{
			if p.Color() == WHITE {
				return 'N'
			}
			return 'n'
		}
	case BishopType:
		{
			if p.Color() == WHITE {
				return 'B'
			}
			return 'b'
		}
	case RookType:
		{
			if p.Color() == WHITE {
				return 'R'
			}
			return 'r'
		}
	case QueenType:
		{
			if p.Color() == WHITE {
				return 'Q'
			}
			return 'q'
		}
	case KingType:
		{
			if p.Color() == WHITE {
				return 'K'
			}
			return 'k'
		}
	default:
		{
			return 'X'
		}
	}
}
