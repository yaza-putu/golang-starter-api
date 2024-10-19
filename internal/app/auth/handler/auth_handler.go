package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/repository"
	"github.com/yaza-putu/golang-starter-api/internal/app/auth/service"
	"github.com/yaza-putu/golang-starter-api/internal/http/request"
	"github.com/yaza-putu/golang-starter-api/internal/http/response"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
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
	req := entity.Token{}
	b := ctx.Bind(&req)
	if b != nil {
		return ctx.JSON(http.StatusBadRequest, response.Api(
			response.SetMessage(b.Error()),
		))
	}

	// validation form
	res, err := request.Validation(&req)
	logger.New(err, logger.SetType(logger.INFO))

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, res)
	}

	// get ip and device
	device := ctx.Request().Header.Get("User-Agent")
	ip := strings.Split(ctx.Request().RemoteAddr, ":")[0]

	r := a.authService.Login(req.Email, req.Password, ip, device)

	if r.GetCode() == 200 {
		// set coockie device ID
		data := r.GetData().(map[string]string)
		cookie := new(http.Cookie)
		cookie.Name = "DVID"
		cookie.Value = data["device_id"]
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = http.SameSiteStrictMode

		http.SetCookie(ctx.Response().Writer, cookie)
	}

	return ctx.JSON(r.GetCode(), r)
}

func (a *authHandler) Refresh(ctx echo.Context) error {
	// request
	req := entity.RefreshToken{}

	b := ctx.Bind(&req)
	if b != nil {
		return ctx.JSON(http.StatusBadRequest, response.Api(
			response.SetMessage(b.Error()),
		))
	}

	// read device id from coockie
	cookie, err := ctx.Cookie("DVID")

	if err == nil && cookie.Value != "" {
		req.DeviceId = cookie.Value
	}
	// validation form
	res, err := request.Validation(&req)
	logger.New(err, logger.SetType(logger.INFO))

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, res)
	}

	r := a.authService.Refresh(req.DeviceId)

	return ctx.JSON(r.GetCode(), r)
}
