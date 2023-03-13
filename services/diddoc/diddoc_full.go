package diddoc

import (
	"github.com/cheqd/did-resolver/services"
)

type FullDIDDocRequestService struct {
	BaseDidDocRequestService
}

func (dd *FullDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *FullDIDDocRequestService) Prepare(c services.ResolverContext) error {
	return nil
}

func (dd *FullDIDDocRequestService) MakeResponse(c services.ResolverContext) error {
	return nil
}
