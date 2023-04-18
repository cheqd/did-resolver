package diddoc

import (
	"sort"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type DidDocResolveHandler struct {
	queries.BaseQueryHandler
	DidDocHelperHandler
}

func (dd *DidDocResolveHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	metadata := c.QueryParams().Get(types.Metadata)
	// If metadata is set we don't need to resolve the DidDoc
	if metadata != "" {
		return dd.Continue(c, service, response)
	}

	allVersions, err := dd.CastToContent(service, response)
	if err != nil {
		return nil, err
	}
	if len(allVersions) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, dd.IsDereferencing)
	}

	// Get the latest version. If versionId and versionTime handlers were called here, should be only 1 element.
	// If versionId or versionTime was not called, we will return the latest version
	versionId := allVersions[0].VersionId
	filteredResources := allVersions[0].Resources
	
	// Filter in descending order
	sort.Sort(filteredResources)

	result, _err := c.DidDocService.Resolve(service.GetDid(), versionId, service.GetContentType())
	if _err != nil {
		_err.IsDereferencing = dd.IsDereferencing
		return nil, _err
	}

	result.Metadata.Resources = filteredResources

	// Call the next handler
	return dd.Continue(c, service, result)
}
