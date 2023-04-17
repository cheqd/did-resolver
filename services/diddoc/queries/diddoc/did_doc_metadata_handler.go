package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type DidDocMetadataHandler struct {
	queries.BaseQueryHandler
	DidDocHelperHandler
}

func (dd *DidDocMetadataHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	metadata := c.QueryParams().Get(types.Metadata)
	// If metadata is set we don't need to resolve the DidDoc
	if metadata == "" {
		return dd.Continue(c, service, response)
	}

	allVersions, err := dd.CastToContent(service, response)
	if err != nil {
		return nil, err
	}
	if len(allVersions) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, dd.IsDereferencing)
	}

	// Get the latest version. If versionId and versionTime handlers were called, should be only 1 element.
	// If versionId or versionTime was not called, we will return the latest version
	// Cause allVersions are sorted in reverse order the latest version is the first element
	versionId := allVersions[0].VersionId
	filteredResources := allVersions[0].Resources
	result, err := c.DidDocService.GetDIDDocVersionsMetadata(service.GetDid(), versionId, service.GetContentType())
	// Fill the resources
	content := result.ContentStream.(*types.ResolutionDidDocMetadata)
	content.Resources = filteredResources
	result.ContentStream = content

	return dd.Continue(c, service, result)
}