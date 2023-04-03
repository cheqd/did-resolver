package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DIDDocAllVersionMetadataRequestService struct {
	services.BaseRequestService
}

func (dd *DIDDocAllVersionMetadataRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = true
	return nil
}

func (dd *DIDDocAllVersionMetadataRequestService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}

func (dd DIDDocAllVersionMetadataRequestService) Redirect(c services.ResolverContext) error {
	path := types.RESOLVER_PATH + dd.Did + types.DID_VERSIONS_PATH
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *DIDDocAllVersionMetadataRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *DIDDocAllVersionMetadataRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.GetAllDidDocVersionsMetadata(dd.Did, dd.RequestedContentType)
	if err != nil {
		err.IsDereferencing = dd.IsDereferencing
		return err
	}
	dd.Result = result
	return nil
}
