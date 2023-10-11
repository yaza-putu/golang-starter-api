package routes

import (
	"github.com/labstack/echo/v4"
	"yaza/src/app/auth/handler"
)

var authhandler = handler.NewAuthHandler()

func Api(r *echo.Echo) {
	r.GET("/token", authhandler.Create)
}
