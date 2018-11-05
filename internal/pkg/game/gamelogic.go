package game

type Result struct {
	isWhiteWinner bool
}

type Turn string

type GameLogic interface {
	FirstPlayerTurn(Turn) error
	SecondPlayerTurn(Turn) error
	Start() (bool, <-chan Result) // true  - 1 игрок играет за белых. false - 2 игрок играет за белых
}
