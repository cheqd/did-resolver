package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceChecksumHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceChecksumHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceChecksum := service.GetQueryParam(types.ResourceChecksum)
	if resourceChecksum == "" {
		return d.Continue(c, service, response)
	}

	// Cast to just list of resources
	didResolution, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	// Filter the list of metadatas by the resourceCollectionId
	resourceCollectionFiltered := didResolution.Metadata.Resources.FilterByChecksum(resourceChecksum)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	didResolution.Metadata.Resources = resourceCollectionFiltered

	// Call the next handler
	return d.Continue(c, service, didResolution)
}
