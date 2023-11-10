package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yaza-putu/golang-starter-api/src/app/auth/handler"
	cat "github.com/yaza-putu/golang-starter-api/src/app/category/handler"
	gds "github.com/yaza-putu/golang-starter-api/src/app/goods/handler"
	"github.com/yaza-putu/golang-starter-api/src/config"
)

var authhandler = handler.NewAuthHandler()
var categoryHandler = cat.NewCategoryHandler()
var goodsHandler = gds.NewGoodsHandler()

func Api(r *echo.Echo) {
	route := r.Group("api")
	{
		route.POST("/token", authhandler.Create)
		route.PUT("/token", authhandler.Refresh)

		auth := route.Group("")
		{
			auth.Use(middleware.JWT([]byte(config.Key().Token)))

			category := auth.Group("/categories")
			{
				category.GET("", categoryHandler.All)
				category.POST("", categoryHandler.Create)
				category.GET("/:id", categoryHandler.FindById)
				category.PUT("/:id", categoryHandler.Update)
				category.DELETE("/:id", categoryHandler.Delete)
			}

			goods := auth.Group("/goods")
			{
				goods.GET("", goodsHandler.All)
				goods.POST("", goodsHandler.Create)
				goods.GET("/:id", goodsHandler.FindById)
				goods.PUT("/:id", goodsHandler.Update)
				goods.PATCH("/:id", goodsHandler.Stock)
				goods.DELETE("/:id", goodsHandler.Delete)
			}
		}
	}
}
