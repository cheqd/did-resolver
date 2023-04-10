package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceCollectionIdHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceCollectionIdHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceCollectionId := service.GetQueryParam(types.ResourceCollectionId)
	if resourceCollectionId == "" {
		return d.Continue(c, service, response)
	}

	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	// Filter the list of metadatas by the resourceCollectionId
	resourceCollectionFiltered := resourceCollection.Resources.FilterByCollectionId(resourceCollectionId)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	// Call the next handler
	return d.Continue(c, service, d.CastToResult(resourceCollectionFiltered))
}
