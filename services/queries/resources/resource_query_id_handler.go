package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceIdHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceIdHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceId := service.GetQueryParam(types.ResourceId)
	if resourceId == "" {
		return d.Continue(c, service, response)
	}

	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}
	resourceCollectionFiltered := resourceCollection.Resources.GetByResourceId(resourceId)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	// Call the next handler
	return d.Continue(c, service, d.CastToResult(resourceCollectionFiltered))
}
