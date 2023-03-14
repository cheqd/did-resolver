package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DIDDocVersionMetadataRequestService struct {
	services.BaseRequestService
}

func (dd *DIDDocVersionMetadataRequestService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}

func (dd DIDDocVersionMetadataRequestService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.Did)

	path := types.RESOLVER_PATH + migratedDid + types.DID_VERSION_PATH + dd.Version
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *DIDDocVersionMetadataRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *DIDDocVersionMetadataRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.GetDIDDocVersionsMetadata(dd.Did, dd.Version, dd.RequestedContentType)
	if err != nil {
		return err
	}
	dd.Result = result
	return nil
}

func (dd *DIDDocVersionMetadataRequestService) MakeResponse(c services.ResolverContext) error {
	return nil
}
