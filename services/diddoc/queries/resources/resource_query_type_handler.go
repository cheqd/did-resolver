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
	// Get the name of the first resource
	firstName := resources[0].Name
	// Compare all other names with the first one
	for _, resource := range resources[1:] {
		if resource.Name != firstName {
			return false
		}
	}
	return true
}
