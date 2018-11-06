package chess

import (
	"fmt"
	"sort"
)

type Game struct {
	board *board // P N B R Q K p n b r q k .
	turn  pieceColor
}

func NewGame() *Game {
	return &Game{
		board: newBoard(),
		turn:  white,
	}
}

func (g *Game) Move(uci string) error {
	legalMoves := g.board.legalMoves(g.turn)

	board, exists := legalMoves[uci]
	if !exists {
		return fmt.Errorf("%s is illegal", uci)
	}
	g.board = board

	if g.turn == white {
		g.turn = black
	} else {
		g.turn = white
	}

	return nil
}

func (g *Game) IsCheckmate() bool {
	return g.board.isCheckmate(g.turn)
}

func (g *Game) IsStalemate() bool {
	return g.board.isStalemate(g.turn)
}

func (g *Game) IsGameOver() bool {
	return g.IsCheckmate() || g.IsStalemate()
}

func (g *Game) PrintBoard() {
	g.board.printBoard()
}

func (g *Game) PrintLegalMoves() {
	var moves []string
	for move := range g.board.legalMoves(g.turn) {
		moves = append(moves, move)
	}
	sort.Strings(moves)
	for _, move := range moves {
		fmt.Printf("%s ", move)
	}
	fmt.Println()
}
