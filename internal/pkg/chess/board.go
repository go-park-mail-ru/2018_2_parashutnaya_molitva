package chess

import (
	"fmt"
)

type board struct {
	field [][]piece
}

// board with starting position
func newBoard() *board {
	field := make([][]piece, 8)

	field[0] = []piece{
		newPiece(rookType, white), newPiece(knightType, white), newPiece(bishopType, white), newPiece(queenType, white),
		newPiece(kingType, white), newPiece(bishopType, white), newPiece(knightType, white), newPiece(rookType, white),
	}
	for i := 1; i < 7; i++ {
		field[i] = make([]piece, 8)
	}
	for i := 0; i < 8; i++ {
		field[1][i] = newPiece(pawnType, white)
	}
	for i := 2; i < 6; i++ {
		for j := 0; j < 8; j++ {
			field[i][j] = newPiece(emptyType, none)
		}
	}
	for i := 0; i < 8; i++ {
		field[6][i] = newPiece(pawnType, black)
	}
	field[7] = []piece{
		newPiece(rookType, black), newPiece(knightType, black), newPiece(bishopType, black), newPiece(queenType, black),
		newPiece(kingType, black), newPiece(bishopType, black), newPiece(knightType, black), newPiece(rookType, black),
	}

	return &board{field: field}
}

// board copy
func (b *board) copy() *board {
	duplicateField := make([][]piece, len(b.field))
	for i := range b.field {
		duplicateField[i] = make([]piece, len(b.field[i]))
		copy(duplicateField[i], b.field[i])
	}

	copiedBoard := board{duplicateField}
	return &copiedBoard
}

// moves piece from `from` to `to`
func (b *board) movePiece(from, to *coord) {
	b.field[to.r][to.c] = b.field[from.r][from.c]
	b.field[to.r][to.c].setMoved(true)
	b.field[from.r][from.c] = newPiece(emptyType, none)
}

// returns piece at coord `pos`
func (b *board) pieceAt(pos *coord) *piece {
	if pos.r < 0 || pos.r >= 8 || pos.c < 0 || pos.c >= 8 {
		nonePiece := newPiece(noneType, none)
		return &nonePiece
	}
	return &b.field[pos.r][pos.c]
}

// sets piece `p` at coord `pos`
func (b *board) setPieceAt(pos *coord, p piece) {
	b.field[pos.r][pos.c] = p
}

// prints current board position
func (b *board) printBoard() {
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			fmt.Printf("%c", b.pieceAt(&coord{i, j}).shortName())
		}
		fmt.Println()
	}
}

// returns available moves for current position `color`
func (b *board) legalMoves(color pieceColor) map[string]*board {
	pseudoLegalMoves := b.pseudoLegalMovesWithColor(color, false)
	legalMoves := make(map[string]*board)
	for move, board := range pseudoLegalMoves {
		if !board.isCheck(color) {
			legalMoves[move] = board
		}
	}
	return legalMoves
}

// returns unfiltered (not checking for check) moves for coord `pos`
// `attackOnly` is required to dodge infinite recursion while solving king moves
func (b *board) pseudoLegalMovesAtPos(pos *coord, attackOnly bool) map[string]*board {
	availableMoves := make(map[string]*board)

	piece := b.pieceAt(pos)
	switch piece.getType() {
	case pawnType:
		{
			pawnAvailableMoves := pawnMoves(b, pos, attackOnly)
			for key, val := range pawnAvailableMoves {
				availableMoves[key] = val
			}
		}
	case knightType:
		{
			knightAvailableMoves := knightMoves(b, pos)
			for key, val := range knightAvailableMoves {
				availableMoves[key] = val
			}
		}
	case bishopType:
		{
			bishopAvailableMoves := bishopMoves(b, pos)
			for key, val := range bishopAvailableMoves {
				availableMoves[key] = val
			}
		}
	case rookType:
		{
			rookAvailableMoves := rookMoves(b, pos)
			for key, val := range rookAvailableMoves {
				availableMoves[key] = val
			}
		}
	case queenType:
		{
			queenAvailableMoves := queenMoves(b, pos)
			for key, val := range queenAvailableMoves {
				availableMoves[key] = val
			}
		}
	case kingType:
		{
			kingAvailableMoves := kingMoves(b, pos, attackOnly)
			for key, val := range kingAvailableMoves {
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

// returns unfiltered (not checking for check) moves for `color`
// `attackOnly` is required to dodge infinite recursion while solving king moves
func (b *board) pseudoLegalMovesWithColor(color pieceColor, attackOnly bool) map[string]*board {
	availableMoves := make(map[string]*board)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.pieceAt(&coord{i, j}).getColor() == color {
				movesAtPos := b.pseudoLegalMovesAtPos(&coord{i, j}, attackOnly)
				for key, val := range movesAtPos {
					availableMoves[key] = val
				}
			}
		}
	}

	return availableMoves
}

// removes en passant piece from board
func (b *board) removeEnPassant() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if (b.pieceAt(&coord{i, j}).getType() == enPassantType) {
				b.setPieceAt(&coord{i, j}, newPiece(emptyType, none))
				break
			}
		}
	}
}

// returns true if `color` is in check
func (b *board) isCheck(color pieceColor) bool {
	oppositeColor := white
	if color == white {
		oppositeColor = black
	}
	pseudoMoves := b.pseudoLegalMovesWithColor(oppositeColor, true)

	for _, moveBoard := range pseudoMoves {
		if !moveBoard.kingExists(color) {
			return true
		}
	}

	return false
}

// returns true if king of `color` exists on board
func (b *board) kingExists(color pieceColor) bool {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			p := b.pieceAt(&coord{i, j})
			if p.getType() == kingType && p.getColor() == color {
				return true
			}
		}
	}
	return false
}

// returns true if `color` is checkmated
func (b *board) isCheckmate(color pieceColor) bool {
	return b.isCheck(color) && len(b.legalMoves(color)) == 0
}

// returns true if stalemate
func (b *board) isStalemate(color pieceColor) bool {
	return !b.isCheck(color) && len(b.legalMoves(color)) == 0
}

// return true if only 2 kings left
func (b *board) isInsufficientMaterial() bool {
	pieceCounter := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			p := b.pieceAt(&coord{i, j})
			if p.getType() != emptyType && p.getType() != noneType && p.pieceType != enPassantType {
				pieceCounter++
				if pieceCounter > 2 {
					return false
				}
			}
		}
	}
	return true
}
