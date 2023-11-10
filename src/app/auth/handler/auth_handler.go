package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/src/app/auth/repository"
	"github.com/yaza-putu/golang-starter-api/src/app/auth/service"
	"github.com/yaza-putu/golang-starter-api/src/app/auth/validation"
	"github.com/yaza-putu/golang-starter-api/src/http/request"
	"github.com/yaza-putu/golang-starter-api/src/http/response"
	"github.com/yaza-putu/golang-starter-api/src/logger"
	"net/http"
	"time"
)

type authHandler struct {
	authService service.AuthInterface
}

func NewAuthHandler() *authHandler {
	return &authHandler{
		authService: service.NewAuthService(repository.NewUserRepository(), service.NewToken()),
	}
}

func (a *authHandler) Create(ctx echo.Context) error {
	// request validation & capture data
	req := validation.TokenValidation{}
	b := ctx.Bind(&req)
	if b != nil {
		return ctx.JSON(http.StatusBadRequest, response.Api(
			response.SetCode(400), response.SetMessage(b.Error()),
		))
	}

	// validation form
	res, err := request.Validation(&req)
	logger.New(err, logger.SetType(logger.INFO))

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, res)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	r := a.authService.Login(ctxTimeout, req.Email, req.Password)

	return ctx.JSON(r.Code, r)
}

func (a *authHandler) Refresh(ctx echo.Context) error {
	// request
	req := validation.RefreshTokenValidation{}

	b := ctx.Bind(&req)
	if b != nil {
		return ctx.JSON(http.StatusBadRequest, response.Api(
			response.SetCode(400), response.SetMessage(b.Error()),
		))
	}

	// validation form
	res, err := request.Validation(&req)
	logger.New(err, logger.SetType(logger.INFO))

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, res)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	r := a.authService.Refresh(ctxTimeout, req.Token)

	return ctx.JSON(r.Code, r)
}
