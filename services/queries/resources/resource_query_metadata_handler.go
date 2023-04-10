package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceMetadataHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceMetadataHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceMetadata := service.GetQueryParam(types.ResourceMetadata)

	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	// After all filters and validation only one resource should be left
	resource := resourceCollection.Resources[0]

	if resourceMetadata == "true" {
		dereferencingResult, err := c.ResourceService.DereferenceResourceMetadata(service.GetDid(), resource.ResourceId, service.GetContentType())
		if err != nil {
			return nil, err
		}
		return d.Continue(c, service, dereferencingResult)
	}
	// var dereferenceResult types.ResolutionResultI
	dereferenceResult, _err := c.ResourceService.DereferenceResourceData(service.GetDid(), resource.ResourceId, service.GetContentType())
	// _err = _err.(*types.IdentityError)
	if _err != nil {
		return nil, _err
	}

	// Call the next handler
	return d.Continue(c, service, dereferenceResult)
}
