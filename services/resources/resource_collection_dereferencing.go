package resources

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceCollectionDereferencingService struct {
	services.BaseRequestService
	ResourceId string
}

func (dr *ResourceCollectionDereferencingService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}

func (dr ResourceCollectionDereferencingService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dr.Did)

	path := types.RESOLVER_PATH + migratedDid + types.DID_METADATA
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dr *ResourceCollectionDereferencingService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dr *ResourceCollectionDereferencingService) Query(c services.ResolverContext) error {
	result, err := c.ResourceService.DereferenceCollectionResources(dr.Did, dr.RequestedContentType)
	if err != nil {
		err.IsDereferencing = true
		return err
	}
	dr.Result = result
	return nil
}

func (dd *ResourceCollectionDereferencingService) MakeResponse(c services.ResolverContext) error {
	return nil
}
