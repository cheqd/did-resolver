package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceTypeHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceTypeHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceType := service.GetQueryParam(types.ResourceType)
	if resourceType == "" {
		return d.Continue(c, service, response)
	}

	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}
	// Filter the list of metadata by the resourceCollectionId
	resourceCollectionFiltered := resourceCollection.Resources.FilterByResourceType(resourceType)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	resourceCollection.Resources = resourceCollectionFiltered

	// Call the next handler
	return d.Continue(c, service, resourceCollection)
}
