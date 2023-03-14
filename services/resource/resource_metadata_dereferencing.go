package resources

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceMetadataDereferencingService struct {
	services.BaseRequestService
	ResourceId string
}

func (dr *ResourceMetadataDereferencingService) IsDereferencing() bool {
	return true
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
	return nil
}

func (dr *ResourceMetadataDereferencingService) Query(c services.ResolverContext) error {
	result, err := c.ResourceService.DereferenceResourceMetadata(dr.ResourceId, dr.Did, dr.RequestedContentType)
	if err != nil {
		err.IsDereferencing = dr.IsDereferencing()
		return err
	}
	dr.Result = result
	return nil
}

func (dd *ResourceMetadataDereferencingService) MakeResponse(c services.ResolverContext) error {
	return nil
}
