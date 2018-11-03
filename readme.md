# golang chess engine

## example
```go
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
```

## output
```
rnbqkbnr
pppppppp
........
........
........
........
PPPPPPPP
RNBQKBNR
a2a3 a2a4 b1a3 b1c3 b2b3 b2b4 c2c3 c2c4 d2d3 d2d4 e2e3 e2e4 f2f3 f2f4 g1f3 g1h3 g2g3 g2g4 h2h3 h2h4 


```