package services

import (
	"net/url"
	"strconv"
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

func GetHighestPriorityHeaderAndParams(acceptHeader string) (string, map[string]string) {
	bestHeader := ""
	bestParams := make(map[string]string)
	highestQ := -1.0
	position := 0

	// Split by ',' to separate multiple headers
	headers := strings.Split(acceptHeader, ",")
	for index, entry := range headers {
		entry = strings.TrimSpace(entry)
		parts := strings.Split(entry, ";")

		params := make(map[string]string)
		q := 1.0

		// Parse parameters
		for _, param := range parts[1:] {
			param = strings.TrimSpace(param)
			keyValue := strings.SplitN(param, "=", 2)
			if len(keyValue) == 2 {
				key, value := keyValue[0], strings.Trim(keyValue[1], "\"")
				params[key] = value
			}
		}

		// Extract q value if present
		if qStr, exists := params["q"]; exists {
			if parsedQ, err := strconv.ParseFloat(qStr, 64); err == nil {
				q = parsedQ
			}
		}

		// Determine the highest priority header
		if q > highestQ || (q == highestQ && position > index) {
			highestQ = q
			position = index
			bestHeader = parts[0] // Only header
			delete(params, "q")   // Remove q from params
			bestParams = params   // Update best parameters
		}
	}

	return bestHeader, bestParams
}
