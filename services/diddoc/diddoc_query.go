package diddoc

import (
	"net/url"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type QueryDIDDocRequestService struct {
	BaseDidDocRequestService
}

func (dd *QueryDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *QueryDIDDocRequestService) Prepare(c services.ResolverContext) error {
	queryRaw, flag := services.PrepareQueries(c)
	queries, err := url.ParseQuery(queryRaw)
	if err != nil {
		return err
	}
	if flag != nil {
		return types.NewRepresentationNotSupportedError(dd.did, dd.requestedContentType, nil, true)
	}
	dd.queries = queries

	version := queries.Get("versionId")
	if version != "" {
		dd.version = version
	}
	return nil
}

func (dd *QueryDIDDocRequestService) MakeResponse(c services.ResolverContext) error {
	return nil
}
