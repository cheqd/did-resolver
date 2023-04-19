package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceVersionTimeHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceVersionTimeHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceVersionTime := service.GetQueryParam(types.ResourceVersionTime)
	if resourceVersionTime == "" {
		return d.Continue(c, service, response)
	}

	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}
	// Get resourceId of the resource with the closest time to the requested time
	resourceList, err := resourceCollection.Resources.FindAllBeforeTime(resourceVersionTime)
	if err != nil {
		return nil, types.NewRepresentationNotSupportedError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}
	if len(resourceList) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	resourceCollection.Resources = resourceList

	// Call the next handler
	return d.Continue(c, service, resourceCollection)
}
