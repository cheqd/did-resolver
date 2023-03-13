package diddoc

import (
	"strings"

	"github.com/cheqd/did-resolver/services"
)

type FragmentDIDDocRequestService struct {
	BaseDidDocRequestService
}

func (dd *FragmentDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *FragmentDIDDocRequestService) Prepare(c services.ResolverContext) error {
	splitted := strings.Split(c.Param("did"), "#")

	if len(splitted) == 2 {
		dd.fragment = splitted[1]
	}
	return nil
}

func (dd *FragmentDIDDocRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.DereferenceSecondary(dd.did, dd.version, dd.fragment, dd.requestedContentType)
	if err != nil {
		err.IsDereferencing = true
	}
	dd.result = result
	return nil
}

func (dd *FragmentDIDDocRequestService) MakeResponse(c services.ResolverContext) error {
	return nil
}
