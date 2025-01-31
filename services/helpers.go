package services

import (
	"net/url"
	"strings"

	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
)

func GetContentType(accept string) types.ContentType {
	// It returns supported ContentType or "" otherwise
	typeList := strings.Split(accept, ",")
	for _, cType := range typeList {
		result := types.ContentType(strings.Split(cType, ";")[0])
		if result == "*/*" || result == types.JSONLD {
			return types.DIDJSONLD
		}
		// Make this place more clearly
		if result.IsSupported() {
			return result
		}
	}

	return ""
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

func ExtractMediaTypeParams(mediaType string) map[string]string {
	// Initialize an empty map to hold the parameters
	params := make(map[string]string)

	// Split by ';' to separate media type from parameters
	parts := strings.Split(mediaType, ";")
	if len(parts) <= 1 {
		return params // No parameters present
	}

	// Iterate over the parts (skipping the first one, which is the media type)
	for _, param := range parts[1:] {
		param = strings.TrimSpace(param) // Trim spaces around the parameter
		// Split parameter by '=' to get the key and value
		keyValue := strings.Split(param, "=")
		if len(keyValue) == 2 {
			// Store the key and value in the map
			params[keyValue[0]] = strings.Trim(keyValue[1], "\"") // Remove quotes from value
		}
	}

	return params
}
