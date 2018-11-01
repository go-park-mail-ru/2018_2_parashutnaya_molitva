package chess

import "fmt"

func PawnMoves(b *Board, pos Coord) map[string]*Board {
	availableMoves := make(map[string]*Board)

	pawn := b.PieceAt(pos).(*Pawn)

	// relative coordinates
	forward := Coord{1, 0}
	doubleForward := Coord{2, 0}
	leftCapture := Coord{1, -1}
	rightCapture := Coord{1, 1}

	// reverse steps for black
	if pawn.Color() == BLACK {
		reverseFactor := &Coord{-1, 1}

		forward = forward.multiply(reverseFactor)
		doubleForward = doubleForward.multiply(reverseFactor)
		leftCapture = leftCapture.multiply(reverseFactor)
		rightCapture = rightCapture.multiply(reverseFactor)
	}

	// field coordinates
	forward = forward.add(&pos)
	doubleForward = doubleForward.add(&pos)
	leftCapture = leftCapture.add(&pos)
	rightCapture = rightCapture.add(&pos)

	pieceAtForward := b.PieceAt(forward)
	if pieceAtForward.Type() == EmptyType {
		moveBoard := b.Copy()
		//b.Move(pos, forward)
		availableMoves[CoordsToUci(pos, forward)] = moveBoard
	}

	fmt.Printf("for pos %d %d\n", pos.r, pos.c)
	for key, _ := range availableMoves {
		fmt.Printf("%s ", key)
	}
	fmt.Println()

	return availableMoves
}
