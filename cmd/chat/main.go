package main

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/app/chat"
	"log"
)

func main() {
	err := chat.StartChatServer()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}