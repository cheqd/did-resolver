package resources

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
)

type ResourceDataDereferencingService struct {
	services.BaseRequestService
	ResourceId string
}

func (dr *ResourceDataDereferencingService) Setup(c services.ResolverContext) error {
	dr.IsDereferencing = true
	return nil
}

func (dr *ResourceDataDereferencingService) SpecificPrepare(c services.ResolverContext) error {
	dr.ResourceId = c.Param("resource")
	return nil
}

func (dr ResourceDataDereferencingService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dr.Did)

	path := types.RESOLVER_PATH + migratedDid + types.RESOURCE_PATH + dr.ResourceId
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dr *ResourceDataDereferencingService) SpecificValidation(c services.ResolverContext) error {
	if !utils.IsValidUUID(dr.ResourceId) {
		return types.NewInvalidDidUrlError(dr.ResourceId, dr.RequestedContentType, nil, dr.IsDereferencing)
	}
	return nil
}

func (dr *ResourceDataDereferencingService) Query(c services.ResolverContext) error {
	result, err := c.ResourceService.DereferenceResourceData(dr.Did, dr.ResourceId, dr.RequestedContentType)
	if err != nil {
		err.IsDereferencing = dr.IsDereferencing
		return err
	}
	dr.Result = result
	return nil
}

func (dr ResourceDataDereferencingService) Respond(c services.ResolverContext) error {
	c.Response().Header().Set(echo.HeaderContentType, dr.Result.GetContentType())

	return c.Blob(http.StatusOK, dr.Result.GetContentType(), dr.Result.GetBytes())
}
