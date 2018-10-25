package game

type Cell struct {
	x byte
	y byte
}

type Turn struct {
	from Cell
	to   Cell
}

type GameLogic interface {
	FirstPlayerTurn(Turn) error
	SecondPlayerTurn(Turn) error
}
