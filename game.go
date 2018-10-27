package chess

type Game struct {
	Board Board // P N B R Q K p n b r q k .
}

func NewGame() Game {
	game := Game{}
	game.Board = NewBoard()
	return game
}
