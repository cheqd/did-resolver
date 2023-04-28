package diddoc

import (
	"strings"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type FragmentDIDDocRequestService struct {
	services.BaseRequestService
}

func (dd *FragmentDIDDocRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = true
	return nil
}

func (dd *FragmentDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	// We not allow query here
	if len(dd.Queries) != 0 {
		return types.NewInvalidDidUrlError(dd.GetDid(), dd.RequestedContentType, nil, dd.IsDereferencing)
	}
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
	result, err := c.DidDocService.DereferenceSecondary(dd.GetDid(), dd.Version, dd.Fragment, dd.GetContentType())
	if err != nil {
		err.IsDereferencing = dd.IsDereferencing
		return err
	}
	return dd.SetResponse(result)
}
