package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
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
	didResolution, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	// Filter the list of metadata by the resourceCollectionId
	resourceCollectionFiltered := didResolution.Metadata.Resources.FilterByCollectionId(resourceCollectionId)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	didResolution.Metadata.Resources = resourceCollectionFiltered

	// Call the next handler
	return d.Continue(c, service, didResolution)
}
