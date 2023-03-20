package cmd

import (
	"github.com/cheqd/did-resolver/utils"
)

// If gzip is not accepted by the client, skip the middleware
func GzipSkipper(c echo.Context) bool {
	return !utils.IsGzipAccepted(c)
}
