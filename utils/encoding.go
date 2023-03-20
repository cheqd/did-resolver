package utils

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/multiformats/go-multibase"
)

func ValidateMultibaseEncoding(data string, expectedEncoding multibase.Encoding) error {
	actualEncoding, _, err := multibase.Decode(data)
	if err != nil {
		return err
	}

	if actualEncoding != expectedEncoding {
		return fmt.Errorf("invalid actualEncoding. expected: %s actual: %s",
			multibase.EncodingToStr[expectedEncoding], multibase.EncodingToStr[actualEncoding])
	}

	return nil
}

func ValidateBase58(data string) error {
	return ValidateMultibaseEncoding(string(multibase.Base58BTC)+data, multibase.Base58BTC)
}

func IsValidBase58(data string) bool {
	return ValidateBase58(data) == nil
}

// Headers Encoding
func IsGzipAccepted(c echo.Context) bool {
	acceptEncoding := c.Request().Header.Get(echo.HeaderAcceptEncoding)
	return strings.Contains(acceptEncoding, "gzip")
}
