package main

import (
	"fmt"
	chess "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/app/chess"
)

func main() {
	game := chess.NewGame()
	for game.Status() == chess.InProgress {
		game.PrintBoard()
		game.PrintLegalMoves()
		var uci string
		fmt.Scanln(&uci)
		err := game.Move(uci)
		if err != nil {
			fmt.Println(err)
		}
	}

	switch game.Status() {
	case chess.WhiteWon:
		fmt.Println("white won")
	case chess.BlackWon:
		fmt.Println("black won")
	case chess.Draw:
		fmt.Println("draw")
	}
}
