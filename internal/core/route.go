package core

import (
	"context"
	"fmt"
	"io"
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
	// i18n
	e.Use(middleware.I18nMiddleware)
	// register api route
	routes.Api(e)

	if config.App().Status == "test" {
		e.HideBanner = true
		e.Logger.(*log.Logger).SetOutput(io.Discard)
	} else {
		// app name
		fmt.Printf("APP NAME : %s ", config.App().Name)
	}

	// gracefully shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

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
