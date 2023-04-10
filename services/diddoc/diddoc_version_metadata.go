package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

type DIDDocVersionMetadataRequestService struct {
	services.BaseRequestService
}

func (dd *DIDDocVersionMetadataRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = true
	return nil
}

func (dd *DIDDocVersionMetadataRequestService) SpecificPrepare(c services.ResolverContext) error {
	// Get Version
	dd.Version = c.Param("version")
	return nil
}

func (dd DIDDocVersionMetadataRequestService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.Did)

	path := types.RESOLVER_PATH + migratedDid + types.DID_VERSION_PATH + dd.Version + types.DID_METADATA
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *DIDDocVersionMetadataRequestService) SpecificValidation(c services.ResolverContext) error {
	if !utils.IsValidUUID(dd.Version) {
		return types.NewInvalidDidUrlError(dd.Version, dd.RequestedContentType, nil, dd.IsDereferencing)
	}
	return nil
}

func (dd *DIDDocVersionMetadataRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.GetDIDDocVersionsMetadata(dd.Did, dd.Version, dd.RequestedContentType)
	if err != nil {
		err.IsDereferencing = dd.IsDereferencing
		return err
	}
	dd.Result = result
	return nil
}
