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
	migratedDid := migrations.MigrateDID(dr.Did)

	path := types.RESOLVER_PATH + migratedDid + types.RESOURCE_PATH + dr.ResourceId + "/metadata"
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dr *ResourceMetadataDereferencingService) SpecificValidation(c services.ResolverContext) error {
	if !utils.IsValidUUID(dr.ResourceId) {
		return types.NewInvalidDIDUrlError(dr.ResourceId, dr.RequestedContentType, nil, dr.IsDereferencing)
	}
	return nil
}

func (dr *ResourceMetadataDereferencingService) Query(c services.ResolverContext) error {
	result, err := c.ResourceService.DereferenceResourceMetadata(dr.Did, dr.ResourceId, dr.RequestedContentType)
	if err != nil {
		err.IsDereferencing = dr.IsDereferencing
		return err
	}
	dr.SetResponse(result)
	return nil
}
