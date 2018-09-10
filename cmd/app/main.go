package main

import (
	"flag"
	"fmt"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/app/auth"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/app/server"
)

var (
	portFlag = flag.Int("port", -1, "Порт на котором запуститься сервер")
)

func main() {
	flag.Parse()
	errChan := make(chan error, 2)
	go server.StartApp(*portFlag, errChan)
	go auth.StartAuth(errChan)

	for {
		select {
		case err := <-errChan:
			fmt.Println(err.Error())
		}
	}

}
