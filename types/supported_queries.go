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

func (s *SupportedQueriesT) IntersectWithUrlValues(values url.Values) []string {
	var result []string
	for k := range values {
		if utils.Contains(*s, k) {
			result = append(result, k)
		}
	}
	return result
}

func (s *SupportedQueriesT) Plus(s2 SupportedQueriesT) SupportedQueriesT {
	var result SupportedQueriesT
	result = append(result, *s...)
	result = append(result, s2...)
	return result
}

var DidSupportedQueries = SupportedQueriesT{
	VersionId,
	VersionTime,
	TransformKey,
	ServiceQ,
	RelativeRef,
	Metadata,
}

var DidResolutionQueries = SupportedQueriesT{
	VersionId,
	VersionTime,
	TransformKey,
}

var ResourceSupportedQueries = SupportedQueriesT{
	ResourceId,
	ResourceCollectionId,
	ResourceName,
	ResourceMetadata,
	ResourceType,
	ResourceVersion,
	ResourceVersionTime,
	ResourceChecksum,
}

var AllSupportedQueries = DidSupportedQueries.Plus(ResourceSupportedQueries)

var SupportedQueriesWithTransformKey = []string{
	VersionId,
	VersionTime,
	ServiceQ,
	RelativeRef,
}

func IsSupportedWithCombinationTransformKeyQuery(values url.Values) bool {
	for query := range values {
		if query == TransformKey {
			continue
		}

		if !utils.Contains(SupportedQueriesWithTransformKey, query) {
			return false
		}
	}

	return true
}
