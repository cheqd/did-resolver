package utils

import (
	"github.com/labstack/echo/v4"
)

// If gzip is not accepted by the client, skip the middleware
func GzipSkipper(c echo.Context) bool {
	return !IsGzipAccepted(c)
}
