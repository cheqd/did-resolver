package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ResourceMetadataHandler struct {
	queries.BaseQueryHandler
	ResourceHelperHandler
}

func (d *ResourceMetadataHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceMetadata := service.GetQueryParam(types.ResourceMetadata)
	// Cast to just list of resources
	didResolution, err := d.CastToContent(service, response)
	if err != nil {
		return nil, err
	}
	acceptHeader := c.Request().Header.Get(echo.HeaderAccept)
	contentType, profile := services.GetPriorityContentType(acceptHeader, d.IsDereferencing)

	// special case for single query and resourceMetadata in query, then IsDereferencing is false
	isSingleQuery := len(c.Request().URL.Query()) == 1 && c.QueryParam(types.ResourceMetadata) != ""
	log.Info().Msgf("IsSingleQuery: %v", isSingleQuery)
	log.Info().Msgf("profile: %v", profile)
	if isSingleQuery && profile != types.W3IDDIDURL {
		d.IsDereferencing = false
	}
	// return didResolution result if dereferencing is false
	if !d.IsDereferencing {
		if resourceMetadata == "false" {
			didResolution.Metadata.Resources = nil
		}
		return d.Continue(c, service, didResolution)
	}

	if resourceMetadata == "true" {
		dereferencingResult := types.NewResourceDereferencingFromContent(
			service.GetDid(), service.GetContentType(), &types.ResolutionDidDocMetadata{Resources: didResolution.Metadata.Resources},
		)
		return d.Continue(c, service, dereferencingResult)
	}

	// If it's not a metadata query let's just get the latest Resource.
	// They are sorted in descending order by default
	resource := didResolution.Metadata.Resources[0]

	if contentType == types.JSONLD && profile == types.W3IDDIDURL {
		dereferenceResult, _err := c.ResourceService.DereferenceResourceDataWithMetadata(service.GetDid(), resource.ResourceId, service.GetContentType())
		if _err != nil {
			return nil, _err
		}

		// Call the next handler
		return d.Continue(c, service, dereferenceResult)
	}

	dereferenceResult, _err := c.ResourceService.DereferenceResourceData(service.GetDid(), resource.ResourceId, service.GetContentType())
	if _err != nil {
		return nil, _err
	}

	// Call the next handler
	return d.Continue(c, service, dereferenceResult)
}
