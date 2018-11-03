package chess

type piece struct {
	pieceType pieceType
	color     pieceColor
	isMoved   bool
}

func newPiece(pt pieceType, color pieceColor) piece {
	return piece{
		pieceType: pt,
		color:     color,
		isMoved:   false,
	}
}

func (p *piece) getType() pieceType {
	return p.pieceType
}

func (p *piece) getColor() pieceColor {
	return p.color
}

func (p *piece) getIsMoved() bool {
	return p.isMoved
}

func (p *piece) setMoved(bool) {
	p.isMoved = true
}

func (p *piece) shortName() rune {

	switch p.getType() {
	case emptyType:
		{
			return '.'
		}
	case noneType:
		{
			return '!'
		}
	case enPassantType:
		{
			if p.getColor() == white {
				return 'E'
			}
			return 'e'
		}
	case pawnType:
		{
			if p.getColor() == white {
				return 'P'
			}
			return 'p'
		}
	case knightType:
		{
			if p.getColor() == white {
				return 'N'
			}
			return 'n'
		}
	case bishopType:
		{
			if p.getColor() == white {
				return 'B'
			}
			return 'b'
		}
	case rookType:
		{
			if p.getColor() == white {
				return 'R'
			}
			return 'r'
		}
	case queenType:
		{
			if p.getColor() == white {
				return 'Q'
			}
			return 'q'
		}
	case kingType:
		{
			if p.getColor() == white {
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
