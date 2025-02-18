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
		// content := types.ResolutionDidDocMetadata(resolutionResult.Metadata)
		// // Call the next handler
		// return d.Continue(c, service, types.DidDocMetadataList{content})
	}
	// Otherwise just use the result from previous handlers
	// But here we need to cast ContentStream to ResolutionDidDocMetadata
	// in case of ResourceDereferencing response
	casted_did_resolution, ok := response.(*types.DidResolution)
	if !ok {
		casted_resource_deref, ok := response.(*types.ResourceDereferencing)
		if !ok {
			return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
		}
		content, ok := casted_resource_deref.ContentStream.(*types.ResolutionDidDocMetadata)
		if !ok {
			return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
		}
		// Call the next handler
		return d.Continue(c, service, content)
	}
	// Call the next handler
	return d.Continue(c, service, casted_did_resolution.Metadata)
}
