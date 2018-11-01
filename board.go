package chess

import "fmt"

type Board struct {
	field [][]Piece
}

func NewBoard() Board {
	field := make([][]Piece, 8)

	field[0] = []Piece{
		NewPiece(RookType, WHITE), NewPiece(KnightType, WHITE), NewPiece(BishopType, WHITE), NewPiece(QueenType, WHITE),
		NewPiece(KingType, WHITE), NewPiece(BishopType, WHITE), NewPiece(KnightType, WHITE), NewPiece(RookType, WHITE),
	}
	for i := 1; i < 7; i++ {
		field[i] = make([]Piece, 8)
	}
	for i := 0; i < 8; i++ {
		field[1][i] = NewPiece(PawnType, WHITE)
	}
	for i := 2; i < 6; i++ {
		for j := 0; j < 8; j++ {
			field[i][j] = NewPiece(EmptyType, NONE)
		}
	}
	for i := 0; i < 8; i++ {
		field[6][i] = NewPiece(PawnType, BLACK)
	}
	field[7] = []Piece{
		NewPiece(RookType, BLACK), NewPiece(KnightType, BLACK), NewPiece(BishopType, BLACK), NewPiece(QueenType, BLACK),
		NewPiece(KingType, BLACK), NewPiece(BishopType, BLACK), NewPiece(KnightType, BLACK), NewPiece(RookType, BLACK),
	}

	return Board{field: field}
}

func (b *Board) Assign(o *Board) {
	b.field = o.field
}

func (b *Board) Copy() *Board {
	duplicateField := make([][]Piece, len(b.field))
	for i := range b.field {
		duplicateField[i] = make([]Piece, len(b.field[i]))
		copy(duplicateField[i], b.field[i])
	}

	copiedBoard := Board{duplicateField}
	return &copiedBoard
}

func (b *Board) MoveUci(uci string) {
	availableMoves := b.AvailableMoves()
	val, exists := availableMoves[uci]
	if exists == false {
		fmt.Println("move is illegal")
		return
	}
	b.Assign(val)
}

func (b *Board) MovePieceUci(uci string) {
	b.MovePiece(UciToCoords(uci))
}

func (b *Board) MovePiece(from, to Coord) {
	b.field[to.r][to.c] = b.field[from.r][from.c]
	b.field[to.r][to.c].SetMoved(true)
	b.field[from.r][from.c] = NewPiece(EmptyType, NONE)
}

func (b *Board) PieceAt(pos Coord) *Piece {
	if pos.r < 0 || pos.r >= 8 || pos.c < 0 || pos.c >= 8 {
		nonePiece := NewPiece(NoneType, NONE)
		return &nonePiece
	}
	return &b.field[pos.r][pos.c]
}

func (b *Board) SetPieceAt(pos Coord, p Piece) {
	b.field[pos.r][pos.c] = p
}

func (b *Board) PrintBoard() {
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			fmt.Printf("%c", b.PieceAt(Coord{i, j}).ShortName())
		}
		fmt.Println()
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			b.AvailableMovesAtPos(Coord{i, j})
		}
	}
}

func (b *Board) AvailableMovesAtPos(pos Coord) map[string]*Board {
	availableMoves := make(map[string]*Board)

	piece := b.PieceAt(pos)
	switch piece.Type() {
	case PawnType:
		{
			pawnAvailableMoves := PawnMoves(b, pos)
			for key, val := range pawnAvailableMoves {
				availableMoves[key] = val
			}
		}
	default:
		{
			return availableMoves
		}
	}
	return availableMoves
}

func (b *Board) AvailableMoves() map[string]*Board {
	availableMoves := make(map[string]*Board)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			movesAtPos := b.AvailableMovesAtPos(Coord{i, j})
			for key, val := range movesAtPos {
				availableMoves[key] = val
			}
		}
	}

	return availableMoves
}
