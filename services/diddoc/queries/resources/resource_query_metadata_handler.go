package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
)

type ResourceMetadataHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceMetadataHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceMetadata := service.GetQueryParam(types.ResourceMetadata)
	acceptHeader := c.Request().Header.Get(echo.HeaderAccept)
	_, profile := services.GetPriorityContentType(acceptHeader, d.IsDereferencing)
	// Cast to just list of resources
	resourceCollection, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	if resourceMetadata == "true" {
		if profile == types.W3IDDIDRES {
			didResolutionMetadata := types.NewResolutionMetadata(service.GetDid(), service.GetContentType(), "")

			didResolutionResult := types.DidResolution{ResolutionMetadata: didResolutionMetadata, Metadata: resourceCollection}
			return d.Continue(c, service, didResolutionResult)
		}

		dereferencingResult := types.NewResourceDereferencingFromContent(service.GetDid(), service.GetContentType(), resourceCollection)
		return d.Continue(c, service, dereferencingResult)
	}

	if profile == types.W3IDDIDRES {
		didResolutionMetadata := types.NewResolutionMetadata(service.GetDid(), service.GetContentType(), "")

		didResolutionResult := types.DidResolution{ResolutionMetadata: didResolutionMetadata, Metadata: resourceCollection}
		return d.Continue(c, service, didResolutionResult)
	}
	// If it's not a metadata query let's just get the latest Resource.
	// They are sorted in descending order by default
	resource := resourceCollection.Resources[0]
	dereferenceResult, _err := c.ResourceService.DereferenceResourceData(service.GetDid(), resource.ResourceId, service.GetContentType())
	if _err != nil {
		return nil, _err
	}

	// Call the next handler
	return d.Continue(c, service, dereferenceResult)
}
