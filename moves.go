package chess

func promotionMap() map[string]PieceType {
	pMap := map[string]PieceType{
		"q": QueenType,
		"r": RookType,
		"b": BishopType,
		"n": KnightType,
	}
	return pMap
}

func PawnMoves(b *Board, pos Coord, attackOnly bool) map[string]*Board {
	availableMoves := make(map[string]*Board)

	pawn := b.PieceAt(pos)

	// relative coordinates
	forward := Coord{1, 0}
	doubleForward := Coord{2, 0}
	capture := Coord{1, -1}
	enPasant := Coord{-1, 0}
	pawnEnPassantRelative := Coord{-1, 0}

	// reverse steps for black
	if pawn.Color() == BLACK {
		reverseFactor := &Coord{-1, 1}

		forward = forward.multiply(reverseFactor)
		doubleForward = doubleForward.multiply(reverseFactor)
		capture = capture.multiply(reverseFactor)
		enPasant = enPasant.multiply(reverseFactor)
		pawnEnPassantRelative = pawnEnPassantRelative.multiply(reverseFactor)
	}

	// field coordinates
	forward = forward.add(&pos)
	doubleForward = doubleForward.add(&pos)
	enPasant = doubleForward.add(&enPasant)

	if !attackOnly {
		pieceAtForward := b.PieceAt(forward)
		if pieceAtForward.Type() == EmptyType {
			uci := CoordsToUcis(pos, forward)
			// promotion
			if pawn.Color() == WHITE && forward.r == 7 ||
				pawn.Color() == BLACK && forward.r == 0 {
				pMap := promotionMap()
				for pKey, pType := range pMap {
					moveBoard := b.Copy()
					moveBoard.MovePiece(pos, forward)
					moveBoard.SetPieceAt(forward, NewPiece(pType, pawn.Color()))
					moveBoard.RemoveEnPassant()
					availableMoves[uci+pKey] = moveBoard
				}
			} else {
				moveBoard := b.Copy()
				moveBoard.MovePiece(pos, forward)
				moveBoard.RemoveEnPassant()
				availableMoves[uci] = moveBoard
			}
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

	captureMultipliers := []Coord{{1, 1}, {1, -1}}

	for i := 0; i < len(captureMultipliers); i++ {
		captureAbs := capture.multiply(&captureMultipliers[i])
		captureAbs = captureAbs.add(&pos)
		pieceAtCapture := b.PieceAt(captureAbs)

		if pieceAtCapture.Color() != pawn.Color() && pieceAtCapture.Color() != NONE {
			uci := CoordsToUcis(pos, captureAbs)
			// promotion
			if pawn.Color() == WHITE && captureAbs.r == 7 ||
				pawn.Color() == BLACK && captureAbs.r == 0 {
				pMap := promotionMap()
				for pKey, pType := range pMap {
					moveBoard := b.Copy()
					moveBoard.MovePiece(pos, captureAbs)
					moveBoard.SetPieceAt(captureAbs, NewPiece(pType, pawn.Color()))
					moveBoard.RemoveEnPassant()
					availableMoves[uci+pKey] = moveBoard
				}
			} else {
				moveBoard := b.Copy()
				moveBoard.MovePiece(pos, captureAbs)
				moveBoard.RemoveEnPassant()
				if pieceAtCapture.Type() == EnPassantType {
					pawnEnPassantAbsolute := pawnEnPassantRelative.add(&captureAbs)
					moveBoard.SetPieceAt(pawnEnPassantAbsolute, NewPiece(EmptyType, NONE))
				}
				availableMoves[uci] = moveBoard
			}
		}
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
		piece := b.PieceAt(steps[i])
		if piece.Color() != knight.Color() && piece.Color() != NONE ||
			piece.Type() == EnPassantType ||
			piece.Type() == EmptyType {

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
			piece := b.PieceAt(stepAbsolute)
			if piece.Type() == EmptyType ||
				piece.Type() == EnPassantType {
				steps = append(steps, stepAbsolute)
				continue
			}
			if piece.Color() == bishop.Color() ||
				piece.Type() == NoneType {
				break
			}
			if piece.Color() != bishop.Color() {
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
			piece := b.PieceAt(stepAbsolute)
			if piece.Type() == EmptyType ||
				piece.Type() == EnPassantType {
				steps = append(steps, stepAbsolute)
				continue
			}
			if piece.Color() == rook.Color() ||
				piece.Type() == NoneType {
				break
			}
			if piece.Color() != rook.Color() {
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
		piece := b.PieceAt(stepAbsolute)
		if piece.Type() == NoneType {
			continue
		}
		if piece.Type() == EmptyType ||
			piece.Type() == EnPassantType ||
			piece.Color() != king.Color() {

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
