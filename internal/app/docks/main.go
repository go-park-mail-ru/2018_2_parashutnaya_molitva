package main

import (
	"fmt"
	_ "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func main() {
	http.HandleFunc("/docks/", httpSwagger.WrapHandler)
	err := http.ListenAndServe(":9090", nil)
	fmt.Print(err)
}
