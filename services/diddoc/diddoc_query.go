package diddoc

import (
	"net/url"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

type QueryDIDDocRequestService struct {
	services.BaseRequestService
}

func (dd *QueryDIDDocRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = true
	return nil
}

func (dd *QueryDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	_, err := url.QueryUnescape(dd.Did)
	if err != nil {
		return types.NewInvalidDIDUrlError(dd.Did, dd.RequestedContentType, err, dd.IsDereferencing)
	}

	// ToDo make list of supported queries
	// For now we support only versionId
	if dd.Version == "" {
		return types.NewRepresentationNotSupportedError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	// Validate that versionId is UUID
	if !utils.IsValidUUID(dd.Version) {
		return types.NewInvalidDIDUrlError(dd.Version, dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	return nil
}

func (dd *QueryDIDDocRequestService) SpecificPrepare(c services.ResolverContext) error {
	queryRaw, flag := services.PrepareQueries(c)
	queries, err := url.ParseQuery(queryRaw)
	if err != nil {
		return err
	}
	if flag != nil {
		return types.NewRepresentationNotSupportedError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing)
	}
	dd.Queries = queries

	version := queries.Get("versionId")
	if version != "" {
		dd.Version = version
	}
	return nil
}
