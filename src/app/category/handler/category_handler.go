package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/src/app/category/entity"
	"github.com/yaza-putu/golang-starter-api/src/app/category/service"
	"github.com/yaza-putu/golang-starter-api/src/app/category/validation"
	"github.com/yaza-putu/golang-starter-api/src/http/request"
	"github.com/yaza-putu/golang-starter-api/src/http/response"
	"github.com/yaza-putu/golang-starter-api/src/logger"
	redis_client "github.com/yaza-putu/golang-starter-api/src/redis"
	"github.com/yaza-putu/golang-starter-api/src/utils"
	"net/http"
	"time"
)

type categoryHandler struct {
	service service.CategoryInterface
}

type optFunc func(*service.CategoryInterface)

// NewCategoryHandler with optional parsing service
// benefit to use without parsing service on call and can parsing service for make mock test
func NewCategoryHandler(opts ...optFunc) *categoryHandler {
	o := service.NewCategoryService()

	for _, fn := range opts {
		fn(&o)
	}

	return &categoryHandler{
		service: o,
	}
}

func (c *categoryHandler) All(ctx echo.Context) error {
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
	r := c.service.All(timeOutCtx, req.Page, req.Take)
	return ctx.JSON(http.StatusOK, r)
}

func (c *categoryHandler) Create(ctx echo.Context) error {
	// filter and validation
	req := validation.Category{}
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

func (c *categoryHandler) Update(ctx echo.Context) error {
	// filter and validation
	id := ctx.Param("id")
	req := validation.CategoryWithIgnore{}
	req.ID = id

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

func (c *categoryHandler) Delete(ctx echo.Context) error {
	// filter and validation
	id := ctx.Param("id")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	r := c.service.Delete(ctxTimeout, id)

	return ctx.JSON(r.Code, r)
}

func (c *categoryHandler) FindById(ctx echo.Context) error {
	// filter and validation
	id := ctx.Param("id")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	dst := entity.Category{}
	err := redis_client.Get(ctxTimeout, id, &dst)
	if err != nil {
		r := c.service.FindById(ctxTimeout, id)

		return ctx.JSON(r.Code, r)
	} else {
		return ctx.JSON(http.StatusOK, response.Api(response.SetCode(200), response.SetData(dst)))
	}

}
