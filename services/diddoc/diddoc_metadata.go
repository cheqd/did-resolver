package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DIDDocMetadataService struct {
	services.BaseRequestService
	Profile string
}

func (dr *DIDDocMetadataService) Setup(c services.ResolverContext) error {
	dr.IsDereferencing = false
	return nil
}

func (dr *DIDDocMetadataService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}

func (dr DIDDocMetadataService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dr.GetDid())

	path := types.RESOLVER_PATH + migratedDid + types.DID_METADATA
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dr *DIDDocMetadataService) SpecificValidation(c services.ResolverContext) error {
	// We only allow one query parameter
	if len(dr.Queries) > 1 {
		return types.NewInvalidDidUrlError(dr.GetDid(), dr.RequestedContentType, nil, dr.IsDereferencing)
	}
	return nil
}

func (dr *DIDDocMetadataService) Query(c services.ResolverContext) error {
	resolution, err := c.ResourceService.ResolveMetadataResources(dr.GetDid(), dr.GetContentType())
	if err != nil {
		err.IsDereferencing = false
		return err
	}
	return dr.SetResponse(resolution)
}
