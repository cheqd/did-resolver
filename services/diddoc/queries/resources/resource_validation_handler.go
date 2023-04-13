package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type ResourceValidationHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceValidationHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	if len(resourceCollection.Resources) != 1 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	// Call the next handler
	return d.Continue(c, service, response)
}
