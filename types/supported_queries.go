package types

import (
	"net/url"

	"github.com/cheqd/did-resolver/utils"
)

type SupportedQueriesT []string

// DiffWithUrlValues returns values which are in url.Values but not in SupportedQueriesT
func (s *SupportedQueriesT) DiffWithUrlValues(values url.Values) []string {
	var result []string
	for k := range values {
		if !utils.Contains(*s, k) {
			result = append(result, k)
		}
	}
	return result
}

var DidSupportedQueries = SupportedQueriesT{
	VersionId,
	VersionTime,
	ServiceQ,
	RelativeRef,
}
