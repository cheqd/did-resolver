package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
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

	// Filter the list of metadata by the resourceCollectionId
	resourceCollectionFiltered := resourceCollection.Resources.FilterByResourceName(resourceName)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	resourceCollection.Resources = resourceCollectionFiltered

	// Call the next handler
	return d.Continue(c, service, resourceCollection)
}
