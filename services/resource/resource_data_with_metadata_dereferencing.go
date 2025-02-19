package resources

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

type ResourceDataWithMetadataDereferencingService struct {
	services.BaseRequestService
	ResourceId string
}

func (dr *ResourceDataWithMetadataDereferencingService) Setup(c services.ResolverContext) error {
	dr.IsDereferencing = true
	return nil
}

func (dr *ResourceDataWithMetadataDereferencingService) SpecificPrepare(c services.ResolverContext) error {
	dr.ResourceId = c.Param("resource")
	return nil
}

func (dr ResourceDataWithMetadataDereferencingService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dr.GetDid())

	path := types.RESOLVER_PATH + migratedDid + types.RESOURCE_PATH + dr.ResourceId
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dr *ResourceDataWithMetadataDereferencingService) SpecificValidation(c services.ResolverContext) error {
	if !utils.IsValidUUID(dr.ResourceId) {
		return types.NewInvalidDidUrlError(dr.ResourceId, dr.RequestedContentType, nil, dr.IsDereferencing)
	}

	// We not allow query here
	if len(dr.Queries) != 0 {
		return types.NewInvalidDidUrlError(dr.GetDid(), dr.RequestedContentType, nil, dr.IsDereferencing)
	}
	return nil
}

func (dr *ResourceDataWithMetadataDereferencingService) Query(c services.ResolverContext) error {
	result, err := c.ResourceService.DereferenceResourceDataWithMetadata(dr.GetDid(), dr.ResourceId, dr.GetContentType())
	if err != nil {
		err.IsDereferencing = dr.IsDereferencing
		return err
	}

	return dr.SetResponse(result)
}
