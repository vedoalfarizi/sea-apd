package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/williamchang80/sea-apd/routes"
)

func main() {
	e := echo.New()
	routes.InitMainRoutes(e)
	appPort := ":" + os.Getenv("APP_PORT")
	appHost := fmt.Sprintf("http://%s%v", os.Getenv("APP_HOST"), appPort)
	fmt.Println("App is running on " + appHost)
	log.Panic(http.ListenAndServe(appPort, e))
}
