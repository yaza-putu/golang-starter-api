package handler

import (
	"github.com/labstack/echo/v4"
	"yaza/src/app/auth/services"
)

type authHandler struct {
	authService services.AuthInterface
}

func NewAuthHandler() *authHandler {
	return &authHandler{authService: services.NewAuthService()}
}

func (a *authHandler) Create(ctx echo.Context) error {
	r := a.authService.Login(ctx.Request().Context(), "admin@mail.com", "Password1")

	return ctx.JSON(r.Code, r)
}
