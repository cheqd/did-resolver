package services

import (
	"net/url"

	"github.com/cheqd/did-resolver/types"
)

type QueryDIDDocRequestService struct {
	BaseDidDocRequestService
}

func (dd *QueryDIDDocRequestService) SpecificValidation(c ResolverContext) error {
	return nil
}

func (dd *QueryDIDDocRequestService) Prepare(c ResolverContext) error {
	queryRaw, flag := prepareQueries(c)
	queries, err := url.ParseQuery(queryRaw)
	if err != nil {
		return err
	}
	if flag != nil {
		return types.NewRepresentationNotSupportedError(dd.did, dd.requestedContentType, nil, true)
	}
	dd.queries = queries
	return nil
}

func (dd *QueryDIDDocRequestService) MakeAnswer(c ResolverContext) error {
	return nil
}