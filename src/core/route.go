package core

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"yaza/src/config"
	"yaza/src/routes"
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
