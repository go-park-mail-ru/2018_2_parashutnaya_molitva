package main

import (
	"flag"
	"log"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/app/server"
)

var (
	portFlag = flag.Int("port", -1, "Порт на котором запуститься сервер")
)

func main() {
	flag.Parse()
	err := server.StartApp(*portFlag)
	if err != nil {
		log.Fatal(err)
	}

}
