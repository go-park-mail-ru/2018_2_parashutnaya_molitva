package main

import (
	"flag"
	"log"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/app/game"
)

var (
	portFlag = flag.Int("port", -1, "Порт на котором запустится сервер")
)

func main() {
	flag.Parse()
	err := game.StartGameServer(*portFlag)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	return
}
