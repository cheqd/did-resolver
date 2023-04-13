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
		// Call the next handler
		return d.Continue(c, service, resolutionResult)
	}
	// Otherwise, we need to dereference the resource using information from previous handlers
	rp, ok := response.(*types.DidResolution)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}
	resolutionResult, err := c.DidDocService.GetDIDDocVersionsMetadata(rp.Did.Id, rp.Metadata.VersionId, service.GetContentType())
	if err != nil {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), err, d.IsDereferencing)
	}
	// Call the next handler
	return d.Continue(c, service, resolutionResult)
}
