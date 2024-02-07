package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/handler"
)

var authhandler = handler.NewAuthHandler()

func Api(r *echo.Echo) {
	route := r.Group("api")
	{
		v1 := route.Group("/v1")
		{
			v1.POST("/token", authhandler.Create)
			v1.PUT("/token", authhandler.Refresh)
		}
	}
}
