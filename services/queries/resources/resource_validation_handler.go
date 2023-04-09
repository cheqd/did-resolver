package resources

import (
	"errors"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/services/queries"
)

type ResourceValidationHandler struct {
	queries.BaseQueryHandler
}

func (d *ResourceValidationHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// After all filters here should be a single resource.
	// Else it's an error
	resDeref, ok := response.(*types.ResourceDereferencing)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceValidationHandler: response is not of type ResourceDereferencing"), d.IsDereferencing)
	}

	// Cast to ResolutionDidDocMetadata for getting the list of metadata
	resourceCollection, ok := resDeref.ContentStream.(*types.ResolutionDidDocMetadata)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceValidationHandler: ContentStream is not of type ResolutionDidDocMetadata"), d.IsDereferencing)
	}

	if len(resourceCollection.Resources) == 0 {
		return nil, types.NewRepresentationNotSupportedError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}
	
	// Call the next handler
	return d.Continue(c, service, response)
}
