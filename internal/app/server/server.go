package server

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
)

var (
	errNoPort = errors.New("Port wasn't passed")
)

func StartApp(port int, errChan chan<- error) {
	if port == -1 {
		errChan <- errNoPort
	}

	stringPort := ":" + strconv.Itoa(port)

	log.Println(port)
	router := routes.NewRouter(http.DefaultServeMux)
	router.HandleFunc("/api/signin", controllers.SignIn).Method("GET")
	router.HandleFunc("/api/signup", controllers.SignUp).Method("POST")
	errChan <- http.ListenAndServe(stringPort, router)
}
