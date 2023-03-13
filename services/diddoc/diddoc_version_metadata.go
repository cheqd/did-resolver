package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DIDDocVersionMetadataRequestService struct {
	BaseDidDocRequestService
}

func (dd *DIDDocVersionMetadataRequestService) Prepare(c services.ResolverContext) error {
	return nil
}

func (dd DIDDocVersionMetadataRequestService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.did)

	path := types.RESOLVER_PATH + migratedDid + types.DID_VERSION_PATH + dd.version
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *DIDDocVersionMetadataRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *DIDDocVersionMetadataRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.GetDIDDocVersionsMetadata(dd.did, dd.version, dd.requestedContentType)
	if err != nil {
		return err
	}
	dd.result = result
	return nil
}

func (dd *DIDDocVersionMetadataRequestService) MakeResponse(c services.ResolverContext) error {
	return nil
}
