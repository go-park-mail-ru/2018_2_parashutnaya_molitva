package chess

func PawnMoves(b *Board, pos Coord) map[string]*Board {
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

	pieceAtForward := b.PieceAt(forward)
	if pieceAtForward.Type() == EmptyType {
		moveBoard := b.Copy()
		moveBoard.MovePiece(pos, forward)
		moveBoard.RemoveEnPassant()
		availableMoves[CoordsToUci(pos, forward)] = moveBoard
	}

	pieceAtDoubleForward := b.PieceAt(doubleForward)
	if pieceAtDoubleForward.Type() == EmptyType && !pawn.IsMoved() {
		moveBoard := b.Copy()
		moveBoard.MovePiece(pos, doubleForward)
		moveBoard.RemoveEnPassant()
		moveBoard.SetPieceAt(enPasant, NewPiece(EnPassantType, pawn.Color()))
		availableMoves[CoordsToUci(pos, doubleForward)] = moveBoard
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
		availableMoves[CoordsToUci(pos, leftCapture)] = moveBoard
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
		availableMoves[CoordsToUci(pos, rightCapture)] = moveBoard
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

			availableMoves[CoordsToUci(pos, steps[i])] = moveBoard
		}
	}

	return availableMoves
}
