package chess

import (
	"fmt"
	"sort"
	"strings"
)

type Game struct {
	board  *board // P N B R Q K p n b r q k .
	turn   pieceColor
	status GameStatus
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

func (g *Game) IsInsufficientMaterial() bool {
	return g.board.isInsufficientMaterial()
}

func (g *Game) IsGameOver() bool {
	return g.IsCheckmate() || g.IsStalemate() || g.IsInsufficientMaterial()
}

// returns true if white
func (g *Game) CurrentTurn() bool {
	return g.turn == white
}

// example: "RNBQKBNRPPPPPPPP................................pppppppprnbqkbnr"
func (g *Game) BoardString() string {
	resultBuilder := &strings.Builder{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Fprintf(resultBuilder, "%c", g.board.pieceAt(&coord{i, j}).shortName())
		}
	}
	return resultBuilder.String()
}

func (g *Game) Status() GameStatus {
	if g.IsCheckmate() {
		if g.turn == white {
			return BlackWon
		} else {
			return WhiteWon
		}
	}
	if g.IsStalemate() || g.IsInsufficientMaterial() {
		return Draw
	}
	return InProgress
}

func (g *Game) LegalMoves() []string {
	legalMoves := g.board.legalMoves(g.turn)
	legalMovesSlice := make([]string, 0, len(legalMoves))
	for key := range legalMoves {
		legalMovesSlice = append(legalMovesSlice, key)
	}

	sort.Strings(legalMovesSlice)
	return legalMovesSlice
}

func (g *Game) PrintBoard() {
	g.board.printBoard()
}

func (g *Game) PrintLegalMoves() {
	fmt.Println(g.LegalMoves())
}
