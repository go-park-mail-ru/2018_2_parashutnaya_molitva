package chess

import "fmt"

func PawnMoves(b *Board, pos Coord, attackOnly bool) map[string]*Board {
	availableMoves := make(map[string]*Board)

	pawn := b.PieceAt(pos)

	// relative coordinates
	forward := Coord{1, 0}
	doubleForward := Coord{2, 0}
	leftCapture := Coord{1, -1}
	rightCapture := Coord{1, 1}
	enPasant := Coord{-1, 0}
	pawnEnPassantRelative := Coord{-1, 0}

	// reverse steps for black
	if pawn.Color() == BLACK {
		reverseFactor := &Coord{-1, 1}

		forward = forward.multiply(reverseFactor)
		doubleForward = doubleForward.multiply(reverseFactor)
		leftCapture = leftCapture.multiply(reverseFactor)
		rightCapture = rightCapture.multiply(reverseFactor)
		enPasant = enPasant.multiply(reverseFactor)
		pawnEnPassantRelative = pawnEnPassantRelative.multiply(reverseFactor)
	}

	// field coordinates
	forward = forward.add(&pos)
	doubleForward = doubleForward.add(&pos)
	leftCapture = leftCapture.add(&pos)
	rightCapture = rightCapture.add(&pos)
	enPasant = doubleForward.add(&enPasant)

	if !attackOnly {
		pieceAtForward := b.PieceAt(forward)
		if pieceAtForward.Type() == EmptyType {
			moveBoard := b.Copy()
			moveBoard.MovePiece(pos, forward)
			moveBoard.RemoveEnPassant()
			availableMoves[CoordsToUcis(pos, forward)] = moveBoard
		}

		pieceAtDoubleForward := b.PieceAt(doubleForward)
		if pieceAtDoubleForward.Type() == EmptyType && !pawn.IsMoved() {
			moveBoard := b.Copy()
			moveBoard.MovePiece(pos, doubleForward)
			moveBoard.RemoveEnPassant()
			moveBoard.SetPieceAt(enPasant, NewPiece(EnPassantType, pawn.Color()))
			availableMoves[CoordsToUcis(pos, doubleForward)] = moveBoard
		}
	}

	pieceAtLeftCapture := b.PieceAt(leftCapture)
	if pieceAtLeftCapture.Color() != pawn.Color() && pieceAtLeftCapture.Color() != NONE {
		moveBoard := b.Copy()
		moveBoard.MovePiece(pos, leftCapture)
		moveBoard.RemoveEnPassant()
		if pieceAtLeftCapture.Type() == EnPassantType {
			pawnEnPassantAbsolute := pawnEnPassantRelative.add(&leftCapture)
			moveBoard.SetPieceAt(pawnEnPassantAbsolute, NewPiece(EmptyType, NONE))
		}
		availableMoves[CoordsToUcis(pos, leftCapture)] = moveBoard
	}

	pieceAtRightCapture := b.PieceAt(rightCapture)
	if pieceAtRightCapture.Color() != pawn.Color() && pieceAtRightCapture.Color() != NONE {
		moveBoard := b.Copy()
		moveBoard.MovePiece(pos, rightCapture)
		moveBoard.RemoveEnPassant()
		if pieceAtRightCapture.Type() == EnPassantType {
			pawnEnPassantAbsolute := pawnEnPassantRelative.add(&rightCapture)
			moveBoard.SetPieceAt(pawnEnPassantAbsolute, NewPiece(EmptyType, NONE))
		}
		availableMoves[CoordsToUcis(pos, rightCapture)] = moveBoard
	}

	return availableMoves
}

func KnightMoves(b *Board, pos Coord) map[string]*Board {
	availableMoves := make(map[string]*Board)

	knight := b.PieceAt(pos)

	steps := []Coord{
		{-2, 1}, {-1, 2}, {1, 2}, {2, 1},
		{2, -1}, {1, -2}, {-1, -2}, {-2, -1},
	}

	// absolute coords
	for i := 0; i < len(steps); i++ {
		steps[i] = steps[i].add(&pos)
	}

	for i := 0; i < len(steps); i++ {
		if b.PieceAt(steps[i]).Color() != knight.Color() && b.PieceAt(steps[i]).Color() != NONE ||
			b.PieceAt(steps[i]).Type() == EnPassantType ||
			b.PieceAt(steps[i]).Type() == EmptyType {

			moveBoard := b.Copy()
			moveBoard.MovePiece(pos, steps[i])
			moveBoard.RemoveEnPassant()

			availableMoves[CoordsToUcis(pos, steps[i])] = moveBoard
		}
	}

	return availableMoves
}

func BishopMoves(b *Board, pos Coord) map[string]*Board {
	availableMoves := make(map[string]*Board)

	bishop := b.PieceAt(pos)

	steps := make([]Coord, 0, 16)

	rMultipliers := []int{1, 1, -1, -1}
	cMultipliers := []int{1, -1, 1, -1}
	for i := 0; i < len(rMultipliers); i++ {
		for j := 1; j < 8; j++ {
			step := Coord{j * rMultipliers[i], j * cMultipliers[i]}
			stepAbsolute := step.add(&pos)
			if b.PieceAt(stepAbsolute).Type() == EmptyType ||
				b.PieceAt(stepAbsolute).Type() == EnPassantType {
				steps = append(steps, stepAbsolute)
				continue
			}
			if b.PieceAt(stepAbsolute).Color() == bishop.Color() ||
				b.PieceAt(stepAbsolute).Type() == NoneType {
				break
			}
			if b.PieceAt(stepAbsolute).Color() != bishop.Color() {
				steps = append(steps, stepAbsolute)
				break
			}
		}
	}

	for i := 0; i < len(steps); i++ {
		moveBoard := b.Copy()
		moveBoard.MovePiece(pos, steps[i])
		moveBoard.RemoveEnPassant()
		availableMoves[CoordsToUcis(pos, steps[i])] = moveBoard
	}

	return availableMoves
}

func RookMoves(b *Board, pos Coord) map[string]*Board {
	availableMoves := make(map[string]*Board)
	rook := b.PieceAt(pos)

	steps := make([]Coord, 0, 16)

	rMultipliers := []int{1, 0, -1, 0}
	cMultipliers := []int{0, 1, 0, -1}
	for i := 0; i < len(rMultipliers); i++ {
		for j := 1; j < 8; j++ {
			step := Coord{j * rMultipliers[i], j * cMultipliers[i]}
			stepAbsolute := step.add(&pos)
			if b.PieceAt(stepAbsolute).Type() == EmptyType ||
				b.PieceAt(stepAbsolute).Type() == EnPassantType {
				steps = append(steps, stepAbsolute)
				continue
			}
			if b.PieceAt(stepAbsolute).Color() == rook.Color() ||
				b.PieceAt(stepAbsolute).Type() == NoneType {
				break
			}
			if b.PieceAt(stepAbsolute).Color() != rook.Color() {
				steps = append(steps, stepAbsolute)
				break
			}
		}
	}

	for i := 0; i < len(steps); i++ {
		moveBoard := b.Copy()
		moveBoard.MovePiece(pos, steps[i])
		moveBoard.RemoveEnPassant()
		availableMoves[CoordsToUcis(pos, steps[i])] = moveBoard
	}

	return availableMoves
}

func QueenMoves(b *Board, pos Coord) map[string]*Board {
	availableMoves := make(map[string]*Board)

	rookAvailableMoves := RookMoves(b, pos)
	for key, val := range rookAvailableMoves {
		availableMoves[key] = val
	}

	bishopAvailableMoves := BishopMoves(b, pos)
	for key, val := range bishopAvailableMoves {
		availableMoves[key] = val
	}

	return availableMoves
}

func KingMoves(b *Board, pos Coord, attackOnly bool) map[string]*Board {
	availableMoves := make(map[string]*Board)

	king := b.PieceAt(pos)

	steps := []Coord{
		{0, -1}, {1, -1}, {1, 0}, {1, 1},
		{0, 1}, {-1, 1}, {-1, 0}, {-1, -1},
	}

	for i := 0; i < len(steps); i++ {
		stepAbsolute := steps[i].add(&pos)
		if b.PieceAt(stepAbsolute).Type() == NoneType {
			continue
		}
		if b.PieceAt(stepAbsolute).Type() == EmptyType ||
			b.PieceAt(stepAbsolute).Type() == EnPassantType ||
			b.PieceAt(stepAbsolute).Color() != king.Color() {

			moveBoard := b.Copy()
			moveBoard.MovePiece(pos, stepAbsolute)
			moveBoard.RemoveEnPassant()
			availableMoves[CoordsToUcis(pos, stepAbsolute)] = moveBoard
		}
	}

	if !attackOnly {
		// king-side castling
		kingSideRookCoord := Coord{0, 3}
		kingSideRookCoord = kingSideRookCoord.add(&pos)

		kingKMovementCoords := []Coord{{0, 1}, {0, 2}}
		for i := 0; i < len(kingKMovementCoords); i++ {
			kingKMovementCoords[i] = kingKMovementCoords[i].add(&pos)
		}

		kRook := b.PieceAt(kingSideRookCoord)
		if !b.IsCheck(king.Color()) && !king.IsMoved() && kRook.Type() == RookType && !kRook.IsMoved() {
			castlingIsLegal := true
			for i := 0; i < len(kingKMovementCoords); i++ {
				if b.PieceAt(kingKMovementCoords[i]).Type() != EmptyType {
					castlingIsLegal = false
					break
				}
				moveBoard := b.Copy()
				moveBoard.MovePiece(pos, kingKMovementCoords[i])
				moveBoard.RemoveEnPassant()
				if moveBoard.IsCheck(king.Color()) {
					castlingIsLegal = false
					break
				}
			}
			if castlingIsLegal {
				fmt.Println("castle is legal")
				moveBoard := b.Copy()
				moveBoard.MovePiece(pos, kingKMovementCoords[len(kingKMovementCoords)-1])
				moveBoard.MovePiece(kingSideRookCoord, kingKMovementCoords[len(kingKMovementCoords)-2])
				moveBoard.RemoveEnPassant()
				availableMoves[CoordsToUcis(pos, kingKMovementCoords[len(kingKMovementCoords)-1])] = moveBoard
			}
		}

		// queen-side castling
		queenSideRookCoord := Coord{0, -4}
		queenSideRookCoord = queenSideRookCoord.add(&pos)

		kingQMovementCoords := []Coord{{0, -1}, {0, -2}}
		for i := 0; i < len(kingQMovementCoords); i++ {
			kingQMovementCoords[i] = kingQMovementCoords[i].add(&pos)
		}

		rookJumpCoord := Coord{0, -3}
		rookJumpCoord = rookJumpCoord.add(&pos)

		qRook := b.PieceAt(queenSideRookCoord)
		if !b.IsCheck(king.Color()) && !king.IsMoved() && qRook.Type() == RookType && !qRook.IsMoved() &&
			b.PieceAt(rookJumpCoord).Type() == EmptyType {
			castlingIsLegal := true
			for i := 0; i < len(kingQMovementCoords); i++ {
				if b.PieceAt(kingQMovementCoords[i]).Type() != EmptyType {
					castlingIsLegal = false
					break
				}
				moveBoard := b.Copy()
				moveBoard.MovePiece(pos, kingQMovementCoords[i])
				moveBoard.RemoveEnPassant()
				if moveBoard.IsCheck(king.Color()) {
					castlingIsLegal = false
					break
				}
			}
			if castlingIsLegal {
				fmt.Println("castle is legal")
				moveBoard := b.Copy()
				moveBoard.MovePiece(pos, kingQMovementCoords[len(kingQMovementCoords)-1])
				moveBoard.MovePiece(queenSideRookCoord, kingQMovementCoords[len(kingQMovementCoords)-2])
				moveBoard.RemoveEnPassant()
				availableMoves[CoordsToUcis(pos, kingQMovementCoords[len(kingQMovementCoords)-1])] = moveBoard
			}
		}

	}
	return availableMoves
}
