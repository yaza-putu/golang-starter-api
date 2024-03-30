package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/golang-starter-api/internal/http/response"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
)

func PanicMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				// Retrieve the stack trace.
				// only capture 4 error stack
				var pcs [4]uintptr
				n := runtime.Callers(3, pcs[:])

				frames := runtime.CallersFrames(pcs[:n])
				errors := []string{}
				for {
					frame, more := frames.Next()
					// capture error with location file, line code and function to easy debug
					errors = append(errors, fmt.Sprintf("\t%s:%d %s\n", frame.File, frame.Line, frame.Function))
					if !more {
						break
					}
				}
				logger.New(fmt.Errorf("panic recover : %v at %v", err, errors))
				c.JSON(http.StatusInternalServerError, response.Api(
					response.SetCode(http.StatusInternalServerError),
					response.SetMessage("Internal Server Error"),
				))
			}
		}()
		return next(c)
	}
}
