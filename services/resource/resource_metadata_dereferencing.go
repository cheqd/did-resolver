package resources

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

type ResourceMetadataDereferencingService struct {
	services.BaseRequestService
	ResourceId string
}

func (dr *ResourceMetadataDereferencingService) Setup(c services.ResolverContext) error {
	dr.IsDereferencing = true
	return nil
}

func (dr *ResourceMetadataDereferencingService) SpecificPrepare(c services.ResolverContext) error {
	dr.ResourceId = c.Param("resource")
	return nil
}

func (dr ResourceMetadataDereferencingService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dr.GetDid())

	path := types.RESOLVER_PATH + migratedDid + types.RESOURCE_PATH + dr.ResourceId + "/metadata"
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dr *ResourceMetadataDereferencingService) SpecificValidation(c services.ResolverContext) error {
	// Metadata endpoint should be one of the supported types.
	if !dr.RequestedContentType.IsSupported() {
		return types.NewRepresentationNotSupportedError(dr.GetDid(), types.JSON, nil, dr.IsDereferencing)
	}
	
	if !utils.IsValidUUID(dr.ResourceId) {
		return types.NewInvalidDidUrlError(dr.ResourceId, dr.RequestedContentType, nil, dr.IsDereferencing)
	}

	// We not allow query here
	if len(dr.Queries) != 0 {
		return types.NewInvalidDidUrlError(dr.GetDid(), dr.RequestedContentType, nil, dr.IsDereferencing)
	}
	return nil
}

func (dr *ResourceMetadataDereferencingService) Query(c services.ResolverContext) error {
	result, err := c.ResourceService.DereferenceResourceMetadata(dr.GetDid(), dr.ResourceId, dr.GetContentType())
	if err != nil {
		err.IsDereferencing = dr.IsDereferencing
		return err
	}
	return dr.SetResponse(result)
}
