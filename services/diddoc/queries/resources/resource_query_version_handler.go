package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceVersionHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceVersionHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceVersion := service.GetQueryParam(types.ResourceVersion)
	if resourceVersion == "" {
		return d.Continue(c, service, response)
	}

	// Cast to just list of resources
	didResolution, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	// Filter the list of metadata by the resourceCollectionId
	resourceCollectionFiltered := didResolution.Metadata.Resources.FilterByResourceVersion(resourceVersion)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	didResolution.Metadata.Resources = resourceCollectionFiltered
	// Call the next handler
	return d.Continue(c, service, didResolution)
}
