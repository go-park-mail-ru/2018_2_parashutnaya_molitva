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

func (b *Board) MoveUCI(uci string) {
	b.Move(UciToCoords(uci))
}

func (b *Board) Move(rFrom, cFrom, rTo, cTo int) {
	b.field[rTo][cTo] = b.field[rFrom][cFrom]
	b.field[rTo][cTo].SetMoved(true)
	b.field[rFrom][cFrom] = NewEmpty()
}

func (b *Board) PieceAt(r, c int) PieceInterface {
	if r < 0 || r >= 8 || c < 0 || c >= 8 {
		return NewNone()
	}
	return b.field[r][c]
}

func (b *Board) PrintBoard() {
	fmt.Println(len(b.field))
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			fmt.Printf("%c", b.PieceAt(i, j).ShortName())
		}
		fmt.Println()
	}
}
