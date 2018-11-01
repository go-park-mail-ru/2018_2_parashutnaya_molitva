package chess

import "fmt"

type Board struct {
	field [][]PieceInterface
}

func NewBoard() Board {
	field := make([][]PieceInterface, 8)

	field[0] = []PieceInterface{
		NewRook(WHITE), NewKnight(WHITE), NewBishop(WHITE), NewQueen(WHITE),
		NewKing(WHITE), NewBishop(WHITE), NewKnight(WHITE), NewRook(WHITE)}
	for i := 1; i < 7; i++ {
		field[i] = make([]PieceInterface, 8)
	}
	for i := 0; i < 8; i++ {
		field[1][i] = NewPawn(WHITE)
	}
	for i := 2; i < 6; i++ {
		for j := 0; j < 8; j++ {
			field[i][j] = NewEmpty()
		}
	}
	for i := 0; i < 8; i++ {
		field[6][i] = NewPawn(BLACK)
	}
	field[7] = []PieceInterface{
		NewRook(BLACK), NewKnight(BLACK), NewBishop(BLACK), NewQueen(BLACK),
		NewKing(BLACK), NewBishop(BLACK), NewKnight(BLACK), NewRook(BLACK)}

	return Board{field: field}
}

func (b *Board) Copy() *Board {
	duplicateField := make([][]PieceInterface, len(b.field))
	for i := range b.field {
		duplicateField[i] = make([]PieceInterface, len(b.field[i]))
		copy(duplicateField[i], b.field[i])
	}

	copiedBoard := Board{duplicateField}
	return &copiedBoard
}

func (b *Board) MoveUCI(uci string) {
	b.Move(UciToCoords(uci))
}

func (b *Board) Move(from, to Coord) {
	b.field[to.r][to.c] = b.field[from.r][from.c]
	b.field[to.r][to.c].SetMoved(true)
	b.field[from.r][from.c] = NewEmpty()
}

func (b *Board) PieceAt(pos Coord) PieceInterface {
	if pos.r < 0 || pos.r >= 8 || pos.c < 0 || pos.c >= 8 {
		return NewNone()
	}
	return b.field[pos.r][pos.c]
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
	//for key, _ := range availableMoves {
	//	fmt.Println(key)
	//}
	return availableMoves
}
