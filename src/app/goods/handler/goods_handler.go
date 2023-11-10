package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/entity"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/service"
	"github.com/yaza-putu/golang-starter-api/src/app/goods/validation"
	"github.com/yaza-putu/golang-starter-api/src/http/request"
	"github.com/yaza-putu/golang-starter-api/src/http/response"
	"github.com/yaza-putu/golang-starter-api/src/logger"
	redis_client "github.com/yaza-putu/golang-starter-api/src/redis"
	"github.com/yaza-putu/golang-starter-api/src/utils"
	"net/http"
	"time"
)

type goodsHandler struct {
	service service.GoodInterface
}

type optFunc func(*service.GoodInterface)

// NewGoodsHandler with optional parsing service
// benefit to use without parsing service on call and can parsing service for make mock test
func NewGoodsHandler(opts ...optFunc) *goodsHandler {
	o := service.NewGoods()

	for _, fn := range opts {
		fn(&o)
	}

	return &goodsHandler{
		service: o,
	}
}

func (c *goodsHandler) All(ctx echo.Context) error {
	// validation and capture request
	req := utils.PaginationValidation{}
	err := ctx.Bind(&req)
	if err != nil {
		logger.New(err, logger.SetType(logger.ERROR))
		return ctx.JSON(http.StatusBadRequest, response.Api(response.SetCode(400), response.SetMessage(err)))
	}
	res, err := request.Validation(&req)
	logger.New(err, logger.SetType(logger.INFO))

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, res)
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	dst := utils.Pagination{}
	er := redis_client.Get(timeOutCtx, "goods", &dst)
	if er != nil {
		r := c.service.All(timeOutCtx, req.Page, req.Take)
		return ctx.JSON(http.StatusOK, r)
	} else {
		return ctx.JSON(http.StatusOK, response.Api(response.SetCode(200), response.SetData(dst)))
	}

}

func (c *goodsHandler) Create(ctx echo.Context) error {
	// filter and validation
	req := validation.Goods{}
	err := ctx.Bind(&req)
	if err != nil {
		logger.New(err, logger.SetType(logger.ERROR))
		return ctx.JSON(http.StatusBadRequest, response.Api(response.SetCode(400), response.SetMessage(err)))
	}

	res, err := request.Validation(&req)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, res)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	r := c.service.Create(ctxTimeout, req)

	return ctx.JSON(r.Code, r)
}

func (c *goodsHandler) Update(ctx echo.Context) error {
	// filter and validation
	id := ctx.Param("id")
	req := validation.Goods{}

	err := ctx.Bind(&req)
	if err != nil {
		logger.New(err, logger.SetType(logger.ERROR))
		return ctx.JSON(http.StatusBadRequest, response.Api(response.SetCode(400), response.SetMessage(err)))
	}

	res, err := request.Validation(&req)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, res)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	r := c.service.Update(ctxTimeout, id, req)

	return ctx.JSON(r.Code, r)
}

func (c *goodsHandler) Stock(ctx echo.Context) error {
	// filter and validation
	id := ctx.Param("id")
	req := validation.GoodsStock{}

	err := ctx.Bind(&req)
	if err != nil {
		logger.New(err, logger.SetType(logger.ERROR))
		return ctx.JSON(http.StatusBadRequest, response.Api(response.SetCode(400), response.SetMessage(err)))
	}

	res, err := request.Validation(&req)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, res)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	r := c.service.Stock(ctxTimeout, id, req.Stock)

	return ctx.JSON(r.Code, r)
}

func (c *goodsHandler) Delete(ctx echo.Context) error {
	// filter and validation
	id := ctx.Param("id")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	r := c.service.Delete(ctxTimeout, id)

	return ctx.JSON(r.Code, r)
}

func (c *goodsHandler) FindById(ctx echo.Context) error {
	// filter and validation
	id := ctx.Param("id")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	dst := entity.Goods{}
	err := redis_client.Get(ctxTimeout, id, &dst)
	if err != nil {
		r := c.service.FindById(ctxTimeout, id)
		return ctx.JSON(r.Code, r)
	} else {
		return ctx.JSON(http.StatusOK, response.Api(response.SetCode(200), response.SetData(dst)))
	}
}
