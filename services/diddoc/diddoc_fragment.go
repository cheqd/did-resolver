package diddoc

import (
	"strings"

	"github.com/cheqd/did-resolver/services"
)

type FragmentDIDDocRequestService struct {
	services.BaseRequestService
}

func (dd *FragmentDIDDocRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = true
	return nil
}

func (dd *FragmentDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *FragmentDIDDocRequestService) SpecificPrepare(c services.ResolverContext) error {
	split := strings.Split(c.Param("did"), "#")

	if len(split) == 2 {
		dd.Fragment = split[1]
	}
	return nil
}

func (dd *FragmentDIDDocRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.DereferenceSecondary(dd.Did, dd.Version, dd.Fragment, dd.RequestedContentType)
	if err != nil {
		err.IsDereferencing = dd.IsDereferencing
		return err
	}
	dd.Result = result
	return nil
}
