package services

import (
	"net/url"
	"strings"

	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
	"github.com/timewasted/go-accept-headers"
)

func GetPriorityContentType(acceptHeader string, resource bool) (types.ContentType, string) {
	// Parse the Accept header using the go-accept-headers package
	acceptedTypes := accept.Parse(acceptHeader)
	if len(acceptedTypes) == 0 {
		// default content type
		return types.JSONLD, ""
	}
	for _, at := range acceptedTypes {
		mediaType := types.ContentType(at.Type + "/" + at.Subtype)

		// If the Header contains any media type, return the default content type
		if mediaType == "*/*" {
			if resource { // If request is from Resource Handlers
				return types.JSONLD, ""
			} else { // If request is from DIDDoc Handlers
				return types.JSONLD, types.W3IDDIDRES
			}
		}

		if resource || mediaType.IsSupported() {
			profile := at.Extensions["profile"]
			profile = strings.Trim(profile, "\"") // Remove surrounding quotes if present
			return mediaType, profile
		}
	}
	return types.ContentType(acceptedTypes[0].Type + "/" + acceptedTypes[0].Subtype), ""
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
