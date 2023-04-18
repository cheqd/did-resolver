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
	resourceType := service.GetQueryParam(types.ResourceType)
	
	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	if resourceType != "" {
		// If we have 2 or more resources we need to check resource names.
		// If resource names are the same we need to return the latest.
		// If resource names are different we need to return an error.
		if !resourceCollection.Resources.AreResourceNamesTheSame() {
			return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
		}
		// They are sorted in descending order by default
		resourceCollection.Resources = types.DereferencedResourceList{resourceCollection.Resources[0]}
	}

	if resourceName != "" {
		// If we have 2 or more resources we need to check resource types.
		// If resource types are the same we need to return the latest.
		// If resource types are different we need to return an error.
		if !resourceCollection.Resources.AreResourceTypesTheSame() {
			return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
		}
		// They are sorted in descending order by default
		resourceCollection.Resources = types.DereferencedResourceList{resourceCollection.Resources[0]}
	}

	// Call the next handler
	return d.Continue(c, service, resourceCollection)
}
