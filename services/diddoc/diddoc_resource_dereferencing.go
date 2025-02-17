package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DIDDocResourceDereferencingService struct {
	services.BaseRequestService
	Profile string
}

func (dr *DIDDocResourceDereferencingService) Setup(c services.ResolverContext) error {
	dr.IsDereferencing = true
	return nil
}

func (dr *DIDDocResourceDereferencingService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}

func (dr DIDDocResourceDereferencingService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dr.GetDid())

	path := types.RESOLVER_PATH + migratedDid + types.DID_METADATA
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dr *DIDDocResourceDereferencingService) SpecificValidation(c services.ResolverContext) error {
	// We only allow one query parameter
	if len(dr.Queries) > 1 {
		return types.NewInvalidDidUrlError(dr.GetDid(), dr.RequestedContentType, nil, dr.IsDereferencing)
	}
	return nil
}

func (dr *DIDDocResourceDereferencingService) Query(c services.ResolverContext) error {
	if dr.Profile == types.W3IDDIDRES {
		resolution, err := c.ResourceService.ResolveCollectionResources(dr.GetDid(), dr.GetContentType())
		if err != nil {
			err.IsDereferencing = false
			return err
		}
		return dr.SetResponse(resolution)
	}
	result, err := c.ResourceService.DereferenceCollectionResources(dr.GetDid(), dr.GetContentType())
	if err != nil {
		err.IsDereferencing = dr.IsDereferencing
		return err
	}

	return dr.SetResponse(result)
}
