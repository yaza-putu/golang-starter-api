package core

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/yaza-putu/golang-starter-api/src/config"
	"github.com/yaza-putu/golang-starter-api/src/routes"
	"io/ioutil"
)

// Routes register and server server
func HttpServe() {
	e := echo.New()
	// set debug mode
	e.Debug = config.App().Debug
	// register api route
	routes.Api(e)

	// app name
	fmt.Printf("APP NAME : %s ", config.App().Name)
	// start server
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s%d", ":", config.Host().Port)))
}

func HttpServerTesting() {
	e := echo.New()
	// set debug mode
	e.Debug = config.App().Debug
	// register api route
	routes.Api(e)
	e.HideBanner = true
	e.Logger.(*log.Logger).SetOutput(ioutil.Discard)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s%d", ":", config.Host().Port)))
}
