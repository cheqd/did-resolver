package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DIDDocAllVersionMetadataRequestService struct {
	BaseDidDocRequestService
}

func (dd *DIDDocAllVersionMetadataRequestService) Prepare(c services.ResolverContext) error {
	return nil
}

func (dd DIDDocAllVersionMetadataRequestService) Redirect(c services.ResolverContext) error {
	path := types.RESOLVER_PATH + dd.did + types.DID_VERSIONS_PATH
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *DIDDocAllVersionMetadataRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *DIDDocAllVersionMetadataRequestService) Query(c services.ResolverContext) error {
	result, rErr := c.DidDocService.GetAllDidDocVersionsMetadata(dd.did, dd.requestedContentType)
	if rErr != nil {
		return rErr
	}
	dd.result = result
	return nil
}

func (dd *DIDDocAllVersionMetadataRequestService) MakeResponse(c services.ResolverContext) error {
	return nil
}
