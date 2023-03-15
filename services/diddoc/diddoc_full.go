package diddoc

import (
	"github.com/cheqd/did-resolver/services"
)

type FullDIDDocRequestService struct {
	services.BaseRequestService
}

func (dd *FullDIDDocRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = false
	return nil
}

func (dd *FullDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *FullDIDDocRequestService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}
