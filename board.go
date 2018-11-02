package chess

import (
	"fmt"
	"sort"
)

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
	availableMoves := b.PseudoLegalMoves(false)
	val, exists := availableMoves[uci]
	if exists == false {
		fmt.Println("move is illegal")
		return
	}
	b.Assign(val)
}

func (b *Board) MovePieceUci(uci string) {
	b.MovePiece(UcisToCoords(uci))
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
	fmt.Println(b.IsCheck(WHITE))
	fmt.Println(b.IsCheck(BLACK))
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			fmt.Printf("%c", b.PieceAt(Coord{i, j}).ShortName())
		}
		fmt.Println()
	}
	b.PrintPseudoLegalMoves()
}

func (b *Board) PrintPseudoLegalMoves() {
	var moves []string
	for move := range b.PseudoLegalMoves(false) {
		moves = append(moves, move)
	}
	sort.Strings(moves)
	for _, move := range moves {
		fmt.Printf("%s ", move)
	}
	fmt.Println()
}

func (b *Board) PseudoLegalMovesAtPos(pos Coord, attackOnly bool) map[string]*Board {
	availableMoves := make(map[string]*Board)

	piece := b.PieceAt(pos)
	switch piece.Type() {
	case PawnType:
		{
			pawnAvailableMoves := PawnMoves(b, pos, attackOnly)
			for key, val := range pawnAvailableMoves {
				availableMoves[key] = val
			}
		}
	case KnightType:
		{
			knightAvailableMoves := KnightMoves(b, pos)
			for key, val := range knightAvailableMoves {
				availableMoves[key] = val
			}
		}
	case BishopType:
		{
			bishopAvailableMoves := BishopMoves(b, pos)
			for key, val := range bishopAvailableMoves {
				availableMoves[key] = val
			}
		}
	case RookType:
		{
			rookAvailableMoves := RookMoves(b, pos)
			for key, val := range rookAvailableMoves {
				availableMoves[key] = val
			}
		}
	case QueenType:
		{
			queenAvailableMoves := QueenMoves(b, pos)
			for key, val := range queenAvailableMoves {
				availableMoves[key] = val
			}
		}
	case KingType:
		{
			kingAvailableMoves := KingMoves(b, pos, attackOnly)
			for key, val := range kingAvailableMoves {
				availableMoves[key] = val
			}
		}
	default:
		{
			return availableMoves
		}
	}

	//for key := range availableMoves {
	//	fmt.Printf("%s ", key)
	//}
	//fmt.Println()
	return availableMoves
}

func (b *Board) PseudoLegalMoves(attackOnly bool) map[string]*Board {
	availableMoves := make(map[string]*Board)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			movesAtPos := b.PseudoLegalMovesAtPos(Coord{i, j}, attackOnly)
			for key, val := range movesAtPos {
				availableMoves[key] = val
			}
		}
	}

	return availableMoves
}

func (b *Board) PseudoLegalMovesWithColor(color PieceColor, attackOnly bool) map[string]*Board {
	availableMoves := make(map[string]*Board)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.PieceAt(Coord{i, j}).Color() == color {
				movesAtPos := b.PseudoLegalMovesAtPos(Coord{i, j}, attackOnly)
				for key, val := range movesAtPos {
					availableMoves[key] = val
				}
			}
		}
	}

	return availableMoves
}

func (b *Board) RemoveEnPassant() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if (b.PieceAt(Coord{i, j}).Type() == EnPassantType) {
				b.SetPieceAt(Coord{i, j}, NewPiece(EmptyType, NONE))
			}
		}
	}
}

func (b *Board) IsCheck(color PieceColor) bool {
	oppositeColor := WHITE
	if color == WHITE {
		oppositeColor = BLACK
	}
	pseudoMoves := b.PseudoLegalMovesWithColor(oppositeColor, true)

	for _, moveBoard := range pseudoMoves {
		if !moveBoard.kingExists(color) {
			return true
		}
	}

	return false
}

func (b *Board) kingExists(color PieceColor) bool {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			p := b.PieceAt(Coord{i, j})
			if p.Type() == KingType && p.Color() == color {
				return true
			}
		}
	}
	return false
}
