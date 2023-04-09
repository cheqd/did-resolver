package queries

import (
	"errors"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceTypeHandler struct {
	BaseQueryHandler
}

func (d *ResourceTypeHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceType := service.GetQueryParam(types.ResourceType)
	if resourceType == "" {
		return d.Continue(c, service, response)
	}

		// If response has type of ResourceDefereferencingResult,
	// then we need to check if the resourceCollectionId is the same as the one in the response
	resDeref, ok := response.(*types.ResourceDereferencing)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceTypeHandler: response is not of type ResourceDereferencing"), d.IsDereferencing)
	}

	// Cast to DereferencedResourceListStruct for getting the list of metadatas
	resourceCollection, ok := resDeref.ContentStream.(*types.ResolutionDidDocMetadata)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceTypeHandler: ContentStream is not of type ResolutionDidDocMetadata"), d.IsDereferencing)
	}
	// Filter the list of metadatas by the resourceCollectionId
	resourceCollectionFiltered := resourceCollection.Resources.FilterByResourceType(resourceType)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	resDeref.ContentStream = &types.ResolutionDidDocMetadata{
		Created:  resourceCollection.Created,
		Updated:  resourceCollection.Updated,
		Deactivated: resourceCollection.Deactivated,
		VersionId: resourceCollection.VersionId,
		Resources: resourceCollectionFiltered,
	}
	
	// Call the next handler
	return d.Continue(c, service, response)
}
