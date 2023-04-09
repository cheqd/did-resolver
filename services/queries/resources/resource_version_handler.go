package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceVersionHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceVersionHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceVersion := service.GetQueryParam(types.ResourceVersion)

	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	// Filter the list of metadata by the resourceCollectionId
	resourceCollectionFiltered := resourceCollection.Resources.FilterByVersion(resourceVersion)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	// Call the next handler
	return d.Continue(c, service, d.CastToResult(resourceCollectionFiltered))
}
