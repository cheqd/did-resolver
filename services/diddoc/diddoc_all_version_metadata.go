package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DIDDocAllVersionMetadataRequestService struct {
	services.BaseRequestService
}

func (dd *DIDDocAllVersionMetadataRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = false
	return nil
}

func (dd *DIDDocAllVersionMetadataRequestService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}

func (dd DIDDocAllVersionMetadataRequestService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.GetDid())

	path := types.RESOLVER_PATH + migratedDid + types.DID_VERSIONS_PATH
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *DIDDocAllVersionMetadataRequestService) SpecificValidation(c services.ResolverContext) error {
	// We not allow query here
	if len(dd.Queries) != 0 {
		return types.NewInvalidDidUrlError(dd.GetDid(), dd.RequestedContentType, nil, dd.IsDereferencing)
	}
	return nil
}

func (dd *DIDDocAllVersionMetadataRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.GetAllDidDocVersionsMetadata(dd.GetDid(), dd.GetContentType())
	if err != nil {
		err.IsDereferencing = dd.IsDereferencing
		return err
	}
	return dd.SetResponse(result)
}
