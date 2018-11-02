package main

import (
	"fmt"
	chess "github.com/Chubasik/chess"
)

func main() {
	game := chess.NewGame()
	for i := 0; i < 1000; i++ {
		game.Board.PrintBoard()
		var uci string
		fmt.Scanln(&uci)
		game.Board.MoveUci(uci)
	}
}
