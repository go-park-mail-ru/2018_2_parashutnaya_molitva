package main

import (
	"flag"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/app/server"
	"log"
)

var (
	portFlag = flag.Int("port", -1, "Порт на котором запустится сервер")
)

func main() {
	flag.Parse()
	err := server.StartApp(*portFlag)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
