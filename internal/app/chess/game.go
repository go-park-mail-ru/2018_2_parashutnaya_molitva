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

// game constructor with starting position
func NewGame() *Game {
	return &Game{
		board: newBoard(),
		turn:  white,
	}
}

// makes move `uci`
// returns nil if move is legal
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

// checkmate condition check
func (g *Game) IsCheckmate() bool {
	return g.board.isCheckmate(g.turn)
}

// stalemate condition check
func (g *Game) IsStalemate() bool {
	return g.board.isStalemate(g.turn)
}

// insufficient material condition check
func (g *Game) IsInsufficientMaterial() bool {
	return g.board.isInsufficientMaterial()
}

// game over condition check
func (g *Game) IsGameOver() bool {
	return g.IsCheckmate() || g.IsStalemate() || g.IsInsufficientMaterial()
}

// returns true if white
func (g *Game) CurrentTurn() bool {
	return g.turn == white
}

// returns flattened board
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

// returns board status
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

// returns slice of legal moves
// example: [a2a3 a2a4 b1a3 b1c3 b2b3 b2b4 c2c3 c2c4 d2d3 d2d4 e2e3 e2e4 f2f3 f2f4 g1f3 g1h3 g2g3 g2g4 h2h3 h2h4]
func (g *Game) LegalMoves() []string {
	legalMoves := g.board.legalMoves(g.turn)
	legalMovesSlice := make([]string, 0, len(legalMoves))
	for key := range legalMoves {
		legalMovesSlice = append(legalMovesSlice, key)
	}

	sort.Strings(legalMovesSlice)
	return legalMovesSlice
}

// prints board
func (g *Game) PrintBoard() {
	g.board.printBoard()
}

// prints legal moves
func (g *Game) PrintLegalMoves() {
	fmt.Println(g.LegalMoves())
}
