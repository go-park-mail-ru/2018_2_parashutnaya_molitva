package chess

import (
	"fmt"
	"sort"
)

type Game struct {
	board *Board // P N B R Q K p n b r q k .
	turn  PieceColor
}

func NewGame() Game {
	game := Game{}
	board := NewBoard()
	game.board = &board
	game.turn = WHITE
	return game
}

func (g *Game) Move(uci string) error {
	legalMoves := g.board.LegalMoves(g.turn)

	board, exists := legalMoves[uci]
	if !exists {
		return fmt.Errorf("%s is illegal", uci)
	}
	g.board = board

	if g.turn == WHITE {
		g.turn = BLACK
	} else {
		g.turn = WHITE
	}

	return nil
}

func (g *Game) IsCheckmate() bool {
	return g.board.isCheckmate(g.turn)
}

func (g *Game) IsStalemate() bool {
	return g.board.isStalemate(g.turn)
}

func (g *Game) PrintBoard() {
	g.board.PrintBoard()
}

func (g *Game) PrintLegalMoves() {
	var moves []string
	for move := range g.board.LegalMoves(g.turn) {
		moves = append(moves, move)
	}
	sort.Strings(moves)
	for _, move := range moves {
		fmt.Printf("%s ", move)
	}
	fmt.Println()
}
