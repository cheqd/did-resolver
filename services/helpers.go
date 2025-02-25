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
	var wildcardFound bool
	var highestPriorityType types.ContentType
	var profile string

	for _, at := range acceptedTypes {
		mediaType, localProfile := extractMediaTypeAndProfile(at)

		// Keep track of the highest priority type in case nothing else matches
		if highestPriorityType == "" {
			highestPriorityType, profile = mediaType, localProfile
		}
		// Detect wildcard "*/*"
		if mediaType == "*/*" {
			wildcardFound = true
			continue // Continue checking other types before making a decision
		}

		// for non-resource query, Check if the media type is supported
		if !resource && mediaType.IsSupported() {
			if profile == "" {
				profile = types.W3IDDIDRES
			}
			return mediaType, profile
		}
	}
	// If the Header contains any media type, return the default content type
	if wildcardFound {
		if resource && profile != types.W3IDDIDURL {
			// If request is from Resource Handlers
			return types.JSONLD, ""
		}
		if !resource { // If request is from DIDDoc Handlers
			return types.JSONLD, types.W3IDDIDRES
		}
	}
	return highestPriorityType, profile
}

// Extracts media type and profile from an accept header entry
func extractMediaTypeAndProfile(at accept.Accept) (types.ContentType, string) {
	mediaType := types.ContentType(at.Type + "/" + at.Subtype)
	profile := strings.Trim(at.Extensions["profile"], "\"") // Remove surrounding quotes if present
	return mediaType, profile
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
