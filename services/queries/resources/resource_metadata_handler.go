package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceMetadataHandler struct {
	queries.BaseQueryHandler
}

func (d *ResourceMetadataHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceMetadata := service.GetQueryParam(types.ResourceMetadata)
	if resourceMetadata == "" {
		return d.Continue(c, service, response)
	}

	// If response has type of ResourceDefereferencingResult,
	// then we need to check if the resourceCollectionId is the same as the one in the response
	// resDeref, ok := response.(*types.ResourceDereferencing)
	// if !ok {
	// 	return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceCollectionIdHandler: response is not of type ResourceDereferencing"), d.IsDereferencing)
	// }

	// // Cast to DereferencedResourceListStruct for getting the list of metadatas
	// resourceCollection, ok := resDeref.ContentStream.(*types.ResolutionDidDocMetadata)
	// if !ok {
	// 	return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceCollectionIdHandler: ContentStream is not of type ResolutionDidDocMetadata"), d.IsDereferencing)
	// }

	// if resourceMetadata != "true" {
	// 	dereferencingResult, err := c.ResourceService.DereferenceResourceMetadata(service.GetDid(), .ResourceId, dr.RequestedContentType)
	// }

	// Call the next handler
	return d.Continue(c, service, response)
}
