package main

import (
	"fmt"
	chess "github.com/Chubasik/chess"
)

func main() {
	game := chess.NewGame()
	for isCheckmate := game.IsCheckmate(); !isCheckmate; isCheckmate = game.IsCheckmate() {
		game.PrintBoard()
		game.PrintLegalMoves()
		var uci string
		fmt.Scanln(&uci)
		err := game.Move(uci)
		if err != nil {
			fmt.Println(err)
		}
	}
}
