package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/i18n"
)

func I18nMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		locale := c.Request().Header.Get("Locale")
		accept := c.Request().Header.Get("Accept-Language")

		if locale != "" {
			i18n.New(i18n.SetLocale(locale))
		} else {
			i18n.New(i18n.SetLocale(accept))
		}

		return next(c)
	}
}
