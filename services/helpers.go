package services

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
	"github.com/timewasted/go-accept-headers"
)

func GetContentType(acceptHeader string) types.ContentType {
	// Parse the Accept header using the go-accept-headers package
	acceptedTypes := accept.Parse(acceptHeader)
	params := make(map[string]string)
	if len(acceptedTypes) == 0 {
		// default content type
		return types.JSONLD
	}

	for _, at := range acceptedTypes {
		mediaType := types.ContentType(at.Type + "/" + at.Subtype)

		if mediaType.IsSupported() {
			fmt.Printf("Selected Media Type: %s, Profile: %s, Q-Value: %f\n", mediaType, at.Extensions["profile"], at.Q)
			return mediaType
		}
		// If the Header contains any media type, return the default content type
		if mediaType == "*/*" {
			params["profile"] = types.W3IDDIDRES
			return types.JSONLD
		}
	}
	return ""
}

func GetContentParams(acceptHeader string) map[string]string {
	// Parse the Accept header using the go-accept-headers package
	acceptedTypes := accept.Parse(acceptHeader)
	params := make(map[string]string)
	if len(acceptedTypes) == 0 {
		// default content type
		return params
	}

	for _, at := range acceptedTypes {
		mediaType := types.ContentType(at.Type + "/" + at.Subtype)

		if mediaType.IsSupported() {
			fmt.Printf("Selected Media Type: %s, Profile: %s, Q-Value: %f\n", mediaType, at.Extensions["profile"], at.Q)
			return at.Extensions
		}
		// If the Header contains any media type, return the default content type
		if mediaType == "*/*" {
			params["profile"] = types.W3IDDIDRES
			return params
		}
	}
	return params
}

func PrepareQueries(c echo.Context) (rawQuery string, flag *string) {
	rawQuery = c.Request().URL.RawQuery
	flagIndex := strings.LastIndex(rawQuery, "%23")
	if flagIndex == -1 || strings.Contains(rawQuery[flagIndex:], "&") {
		return rawQuery, nil
	}
	queryFlag := rawQuery[flagIndex:]

	return rawQuery[0:flagIndex], &queryFlag
}

func GetDidParam(c echo.Context) (string, error) {
	return url.QueryUnescape(c.Param("did"))
}
