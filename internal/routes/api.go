package routes

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/handler"
	"github.com/yaza-putu/golang-starter-api/internal/config"
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

		// auth
		auth := v1.Group("", echojwt.JWT([]byte(config.Key().Token)))
		{
			auth.GET("/whoami", func(c echo.Context) error {
				return c.JSON(http.StatusOK, "OK")
			})
		}
	}
}
