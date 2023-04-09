package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/services/queries"
)

type ResourceNameHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceNameHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceName := service.GetQueryParam(types.ResourceName)

	if resourceName == "" {
		return d.Continue(c, service, response)
	}

	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	// Filter the list of metadatas by the resourceCollectionId
	resourceCollectionFiltered := resourceCollection.Resources.FilterByResourceName(resourceName)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}
	
	// Call the next handler
	return d.Continue(c, service, d.CastToResult(resourceCollectionFiltered))
}
