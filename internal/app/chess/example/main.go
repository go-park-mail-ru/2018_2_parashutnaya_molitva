package main

import (
	"fmt"
	chess "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/app/chess"
)

func main() {
	game := chess.NewGame()
	for !game.IsGameOver() {
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
