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

	didResolution, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}
	// Filter the list of metadata by the resourceCollectionId
	resourceCollectionFiltered := didResolution.Metadata.Resources.FilterByResourceType(resourceType)
	if len(resourceCollectionFiltered) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}
	// If there are multiple resources with the different names or ResourceMetadata=false, return an error
	if len(resourceCollectionFiltered) > 0 && (!HasUniformName(resourceCollectionFiltered) || c.QueryParam(types.ResourceMetadata) == "false") {
		return nil, types.NewInvalidDidUrlError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
	}

	didResolution.Metadata.Resources = resourceCollectionFiltered

	// Call the next handler
	return d.Continue(c, service, didResolution)
}

func HasUniformName(resources types.DereferencedResourceList) bool {
	length := len(resources)
	if length <= 1 {
		return true // Empty or single-element list is considered uniform
	}

	firstName := resources[0].Name
	for i := 1; i < length; i++ {
		if resources[i].Name != firstName {
			return false
		}
	}
	return true
}
