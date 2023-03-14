package diddoc

import (
	"net/url"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type QueryDIDDocRequestService struct {
	services.BaseRequestService
}

func (dd QueryDIDDocRequestService) IsDereferencing() bool {
	return true
}

func (dd *QueryDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *QueryDIDDocRequestService) SpecificPrepare(c services.ResolverContext) error {
	queryRaw, flag := services.PrepareQueries(c)
	queries, err := url.ParseQuery(queryRaw)
	if err != nil {
		return err
	}
	if flag != nil {
		return types.NewRepresentationNotSupportedError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing())
	}
	dd.Queries = queries

	version := queries.Get("versionId")
	if version != "" {
		dd.Version = version
	}
	return nil
}

func (dd *QueryDIDDocRequestService) MakeResponse(c services.ResolverContext) error {
	return nil
}
