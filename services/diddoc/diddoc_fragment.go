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
	splitted := strings.Split(c.Param("did"), "#")

	if len(splitted) == 2 {
		dd.Fragment = splitted[1]
	}
	return nil
}

func (dd *FragmentDIDDocRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.DereferenceSecondary(dd.Did, dd.Version, dd.Fragment, dd.RequestedContentType)
	if err != nil {
		err.IsDereferencing = true
		return err
	}
	dd.Result = result
	return nil
}
