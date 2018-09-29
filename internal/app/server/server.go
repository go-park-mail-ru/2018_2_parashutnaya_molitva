package server

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/controllers"
	//_ "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/docs"
	//httpSwagger "github.com/swaggo/http-swagger"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/routes"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/pkg/errors"
)

var (
	errNoPort = errors.New("Port wasn't passed")
)

func StartApp(port int) error {
	if port == -1 {
		return errNoPort
	}

	stringPort := ":" + strconv.Itoa(port)

	singletoneLogger.LogMessage("Server starting at " + stringPort)
	router := routes.NewRouter(http.DefaultServeMux)
	signIn := middlewareChain(controllers.SignIn, authMiddleware)
	router.HandleFunc("/api/session/", signIn).Method("POST")
	router.HandleFunc("/api/signup", controllers.SignUp).Method("POST")
	router.HandleFunc("/user/:id", controllers.User)
	router.HandleFunc("/user/", controllers.User)

	http.HandleFunc("/docks/",httpSwagger.WrapHandler)
	 err:= http.ListenAndServe(":9090", nil)
	return http.ListenAndServe(stringPort, router)
}
