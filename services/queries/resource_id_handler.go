package queries

import (
	"errors"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceIdHandler struct {
	BaseQueryHandler
}

func (d *ResourceIdHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceId := service.GetQueryParam(types.ResourceId)
	if resourceId == "" {
		return d.Continue(c, service, response)
	}

	// If response has type of ResourceDefereferencingResult,
	// then we need to check if the resourceCollectionId is the same as the one in the response
	resDeref, ok := response.(*types.ResourceDereferencing)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceIdHandler: response is not of type ResourceDereferencing"), d.IsDereferencing)
	}

	resourceCollection, ok := resDeref.ContentStream.(*types.ResolutionDidDocMetadata)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceIdHandler: ContentStream is not of type ResolutionDidDocMetadata"), d.IsDereferencing)
	}
	resourceCollectionFiltered := resourceCollection.Resources.GetByResourceId(resourceId)
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
