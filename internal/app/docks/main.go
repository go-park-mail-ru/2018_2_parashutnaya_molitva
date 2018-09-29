package main

import (
	"net/http"
	_ "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"fmt"
)

func main() {
	http.HandleFunc("/docks/",httpSwagger.WrapHandler)
	err:= http.ListenAndServe(":9090", nil)
	fmt.Print(err)
}
