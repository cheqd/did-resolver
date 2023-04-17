package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceQueryHandler struct {
	queries.BaseQueryHandler
}

func (d *ResourceQueryHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// If response is nil, then we need to dereference the resource from the beginning
	if response == nil {
		resolutionResult, err := c.ResourceService.DereferenceCollectionResources(service.GetDid(), service.GetContentType())
		if err != nil {
			return nil, err
		}
		content, ok := resolutionResult.ContentStream.(*types.ResolutionDidDocMetadata)
		if !ok {
			return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
		}
		// Call the next handler
		return d.Continue(c, service, types.DidDocMetadataList{*content})
	}
	// Otherwise just use the result from previous handlers
	// Call the next handler
	casted, ok := response.(*types.DidResolution)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}
	return d.Continue(c, service, &casted.Metadata)
}
