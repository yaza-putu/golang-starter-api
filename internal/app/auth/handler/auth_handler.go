package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/repository"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/service"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/validation"
	"github.com/yaza-putu/golang-starter-api/internal/http/request"
	"github.com/yaza-putu/golang-starter-api/internal/http/response"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	"net/http"
)

type authHandler struct {
	authService service.Auth
}

func NewAuthHandler() *authHandler {
	return &authHandler{
		authService: service.NewAuth(repository.NewUser(), service.NewToken()),
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

	r := a.authService.Login(ctx.Request().Context(), req.Email, req.Password)

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

	r := a.authService.Refresh(ctx.Request().Context(), req.Token)

	return ctx.JSON(r.Code, r)
}
