package chess

func promotionMap() map[string]pieceType {
	pMap := map[string]pieceType{
		"q": queenType,
		"r": rookType,
		"b": bishopType,
		"n": knightType,
	}
	return pMap
}

func pawnMoves(b *board, pos *coord, attackOnly bool) map[string]*board {
	availableMoves := make(map[string]*board)

	pawn := b.pieceAt(pos)

	// relative coordinates
	forward := &coord{1, 0}
	doubleForward := &coord{2, 0}
	capture := &coord{1, -1}
	enPasant := &coord{-1, 0}
	pawnEnPassantRelative := &coord{-1, 0}

	// reverse steps for black
	if pawn.getColor() == black {
		reverseFactor := &coord{-1, 1}

		forward = forward.multiply(reverseFactor)
		doubleForward = doubleForward.multiply(reverseFactor)
		capture = capture.multiply(reverseFactor)
		enPasant = enPasant.multiply(reverseFactor)
		pawnEnPassantRelative = pawnEnPassantRelative.multiply(reverseFactor)
	}

	// field coordinates
	forward = forward.add(pos)
	doubleForward = doubleForward.add(pos)
	enPasant = doubleForward.add(enPasant)

	if !attackOnly {
		pieceAtForward := b.pieceAt(forward)
		if pieceAtForward.getType() == emptyType {
			uci := mustCoordsToUcis(pos, forward)
			// promotion
			if pawn.getColor() == white && forward.r == 7 ||
				pawn.getColor() == black && forward.r == 0 {
				pMap := promotionMap()
				for pKey, pType := range pMap {
					moveBoard := b.copy()
					moveBoard.movePiece(pos, forward)
					moveBoard.setPieceAt(forward, newPiece(pType, pawn.getColor()))
					moveBoard.removeEnPassant()
					availableMoves[uci+pKey] = moveBoard
				}
			} else {
				moveBoard := b.copy()
				moveBoard.movePiece(pos, forward)
				moveBoard.removeEnPassant()
				availableMoves[uci] = moveBoard
			}
		}

		pieceAtDoubleForward := b.pieceAt(doubleForward)
		if pieceAtDoubleForward.getType() == emptyType && !pawn.getIsMoved() {
			moveBoard := b.copy()
			moveBoard.movePiece(pos, doubleForward)
			moveBoard.removeEnPassant()
			moveBoard.setPieceAt(enPasant, newPiece(enPassantType, pawn.getColor()))
			availableMoves[mustCoordsToUcis(pos, doubleForward)] = moveBoard
		}
	}

	captureMultipliers := []coord{{1, 1}, {1, -1}}

	for i := 0; i < len(captureMultipliers); i++ {
		captureAbs := capture.multiply(&captureMultipliers[i])
		captureAbs = captureAbs.add(pos)
		pieceAtCapture := b.pieceAt(captureAbs)

		if pieceAtCapture.getColor() != pawn.getColor() && pieceAtCapture.getColor() != none {
			uci := mustCoordsToUcis(pos, captureAbs)
			// promotion
			if pawn.getColor() == white && captureAbs.r == 7 ||
				pawn.getColor() == black && captureAbs.r == 0 {
				pMap := promotionMap()
				for pKey, pType := range pMap {
					moveBoard := b.copy()
					moveBoard.movePiece(pos, captureAbs)
					moveBoard.setPieceAt(captureAbs, newPiece(pType, pawn.getColor()))
					moveBoard.removeEnPassant()
					availableMoves[uci+pKey] = moveBoard
				}
			} else {
				moveBoard := b.copy()
				moveBoard.movePiece(pos, captureAbs)
				moveBoard.removeEnPassant()
				if pieceAtCapture.getType() == enPassantType {
					pawnEnPassantAbsolute := pawnEnPassantRelative.add(captureAbs)
					moveBoard.setPieceAt(pawnEnPassantAbsolute, newPiece(emptyType, none))
				}
				availableMoves[uci] = moveBoard
			}
		}
	}

	return availableMoves
}

func knightMoves(b *board, pos *coord) map[string]*board {
	availableMoves := make(map[string]*board)

	knight := b.pieceAt(pos)

	steps := []*coord{
		{-2, 1}, {-1, 2}, {1, 2}, {2, 1},
		{2, -1}, {1, -2}, {-1, -2}, {-2, -1},
	}

	// absolute coords
	for i := 0; i < len(steps); i++ {
		steps[i] = steps[i].add(pos)
	}

	for i := 0; i < len(steps); i++ {
		piece := b.pieceAt(steps[i])
		if piece.getColor() != knight.getColor() && piece.getColor() != none ||
			piece.getType() == enPassantType ||
			piece.getType() == emptyType {

			moveBoard := b.copy()
			moveBoard.movePiece(pos, steps[i])
			moveBoard.removeEnPassant()

			availableMoves[mustCoordsToUcis(pos, steps[i])] = moveBoard
		}
	}

	return availableMoves
}

func bishopMoves(b *board, pos *coord) map[string]*board {
	availableMoves := make(map[string]*board)

	bishop := b.pieceAt(pos)

	steps := make([]*coord, 0, 16)

	rMultipliers := []int{1, 1, -1, -1}
	cMultipliers := []int{1, -1, 1, -1}
	for i := 0; i < len(rMultipliers); i++ {
		for j := 1; j < 8; j++ {
			step := coord{j * rMultipliers[i], j * cMultipliers[i]}
			stepAbsolute := step.add(pos)
			piece := b.pieceAt(stepAbsolute)
			if piece.getType() == emptyType ||
				piece.getType() == enPassantType {
				steps = append(steps, stepAbsolute)
				continue
			}
			if piece.getColor() == bishop.getColor() ||
				piece.getType() == noneType {
				break
			}
			if piece.getColor() != bishop.getColor() {
				steps = append(steps, stepAbsolute)
				break
			}
		}
	}

	for i := 0; i < len(steps); i++ {
		moveBoard := b.copy()
		moveBoard.movePiece(pos, steps[i])
		moveBoard.removeEnPassant()
		availableMoves[mustCoordsToUcis(pos, steps[i])] = moveBoard
	}

	return availableMoves
}

func rookMoves(b *board, pos *coord) map[string]*board {
	availableMoves := make(map[string]*board)
	rook := b.pieceAt(pos)

	steps := make([]*coord, 0, 16)

	rMultipliers := []int{1, 0, -1, 0}
	cMultipliers := []int{0, 1, 0, -1}
	for i := 0; i < len(rMultipliers); i++ {
		for j := 1; j < 8; j++ {
			step := coord{j * rMultipliers[i], j * cMultipliers[i]}
			stepAbsolute := step.add(pos)
			piece := b.pieceAt(stepAbsolute)
			if piece.getType() == emptyType ||
				piece.getType() == enPassantType {
				steps = append(steps, stepAbsolute)
				continue
			}
			if piece.getColor() == rook.getColor() ||
				piece.getType() == noneType {
				break
			}
			if piece.getColor() != rook.getColor() {
				steps = append(steps, stepAbsolute)
				break
			}
		}
	}

	for i := 0; i < len(steps); i++ {
		moveBoard := b.copy()
		moveBoard.movePiece(pos, steps[i])
		moveBoard.removeEnPassant()
		availableMoves[mustCoordsToUcis(pos, steps[i])] = moveBoard
	}

	return availableMoves
}

func queenMoves(b *board, pos *coord) map[string]*board {
	availableMoves := make(map[string]*board)

	rookAvailableMoves := rookMoves(b, pos)
	for key, val := range rookAvailableMoves {
		availableMoves[key] = val
	}

	bishopAvailableMoves := bishopMoves(b, pos)
	for key, val := range bishopAvailableMoves {
		availableMoves[key] = val
	}

	return availableMoves
}

func kingMoves(b *board, pos *coord, attackOnly bool) map[string]*board {
	availableMoves := make(map[string]*board)

	king := b.pieceAt(pos)

	steps := []*coord{
		{0, -1}, {1, -1}, {1, 0}, {1, 1},
		{0, 1}, {-1, 1}, {-1, 0}, {-1, -1},
	}

	for i := 0; i < len(steps); i++ {
		stepAbsolute := steps[i].add(pos)
		piece := b.pieceAt(stepAbsolute)
		if piece.getType() == noneType {
			continue
		}
		if piece.getType() == emptyType ||
			piece.getType() == enPassantType ||
			piece.getColor() != king.getColor() {

			moveBoard := b.copy()
			moveBoard.movePiece(pos, stepAbsolute)
			moveBoard.removeEnPassant()
			availableMoves[mustCoordsToUcis(pos, stepAbsolute)] = moveBoard
		}
	}

	if !attackOnly {
		// king-side castling
		kingSideRookCoord := &coord{0, 3}
		kingSideRookCoord = kingSideRookCoord.add(pos)

		kingKMovementCoords := []*coord{{0, 1}, {0, 2}}
		for i := 0; i < len(kingKMovementCoords); i++ {
			kingKMovementCoords[i] = kingKMovementCoords[i].add(pos)
		}

		kRook := b.pieceAt(kingSideRookCoord)
		if !b.isCheck(king.getColor()) && !king.getIsMoved() && kRook.getType() == rookType && !kRook.getIsMoved() {
			castlingIsLegal := true
			for i := 0; i < len(kingKMovementCoords); i++ {
				if b.pieceAt(kingKMovementCoords[i]).getType() != emptyType {
					castlingIsLegal = false
					break
				}
				moveBoard := b.copy()
				moveBoard.movePiece(pos, kingKMovementCoords[i])
				moveBoard.removeEnPassant()
				if moveBoard.isCheck(king.getColor()) {
					castlingIsLegal = false
					break
				}
			}
			if castlingIsLegal {
				moveBoard := b.copy()
				moveBoard.movePiece(pos, kingKMovementCoords[len(kingKMovementCoords)-1])
				moveBoard.movePiece(kingSideRookCoord, kingKMovementCoords[len(kingKMovementCoords)-2])
				moveBoard.removeEnPassant()
				availableMoves[mustCoordsToUcis(pos, kingKMovementCoords[len(kingKMovementCoords)-1])] = moveBoard
			}
		}

		// queen-side castling
		queenSideRookCoord := &coord{0, -4}
		queenSideRookCoord = queenSideRookCoord.add(pos)

		kingQMovementCoords := []*coord{{0, -1}, {0, -2}}
		for i := 0; i < len(kingQMovementCoords); i++ {
			kingQMovementCoords[i] = kingQMovementCoords[i].add(pos)
		}

		rookJumpCoord := &coord{0, -3}
		rookJumpCoord = rookJumpCoord.add(pos)

		qRook := b.pieceAt(queenSideRookCoord)
		if !b.isCheck(king.getColor()) && !king.getIsMoved() && qRook.getType() == rookType && !qRook.getIsMoved() &&
			b.pieceAt(rookJumpCoord).getType() == emptyType {
			castlingIsLegal := true
			for i := 0; i < len(kingQMovementCoords); i++ {
				if b.pieceAt(kingQMovementCoords[i]).getType() != emptyType {
					castlingIsLegal = false
					break
				}
				moveBoard := b.copy()
				moveBoard.movePiece(pos, kingQMovementCoords[i])
				moveBoard.removeEnPassant()
				if moveBoard.isCheck(king.getColor()) {
					castlingIsLegal = false
					break
				}
			}
			if castlingIsLegal {
				moveBoard := b.copy()
				moveBoard.movePiece(pos, kingQMovementCoords[len(kingQMovementCoords)-1])
				moveBoard.movePiece(queenSideRookCoord, kingQMovementCoords[len(kingQMovementCoords)-2])
				moveBoard.removeEnPassant()
				availableMoves[mustCoordsToUcis(pos, kingQMovementCoords[len(kingQMovementCoords)-1])] = moveBoard
			}
		}

	}
	return availableMoves
}
