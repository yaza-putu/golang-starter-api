package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/http/middleware"
	"github.com/yaza-putu/golang-starter-api/internal/routes"
)

// Routes register and server server
func HttpServe() {
	e := echo.New()
	// set debug mode
	e.Debug = config.App().Debug
	// handle error panic in http req & res
	e.Use(middleware.PanicMiddleware)
	// register api route
	routes.Api(e)

	// gracefully shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// app name
	fmt.Printf("APP NAME : %s ", config.App().Name)

	// start server
	go func() {
		if err := e.Start(fmt.Sprintf("%s%d", ":", config.Host().Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
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
