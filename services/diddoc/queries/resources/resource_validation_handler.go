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
	resourceName := service.GetQueryParam(types.ResourceName)
	resourceMetadata := service.GetQueryParam(types.ResourceMetadata)

	// Cast to just list of resources
	didResolution, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	if resourceMetadata == "true" {
		return d.Continue(c, service, didResolution)
	}

	if resourceName != "" {
		// If we have 2 or more resources we need to check resource types.
		// If resource types are the same we need to return the latest.
		// If resource types are different we need to return an error.
		if !didResolution.Metadata.Resources.AreResourceTypesTheSame() {
			return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
		}
		// They are sorted in descending order by default
		didResolution.Metadata.Resources = types.DereferencedResourceList{didResolution.Metadata.Resources[0]}
	}

	// Call the next handler
	return d.Continue(c, service, didResolution)
}
