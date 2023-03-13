package services

import (
	"github.com/cheqd/did-resolver/types"
)

type FullDIDDocRequestService struct {
	BaseDidDocRequestService
	result *types.DidResolution
}

func (dd *FullDIDDocRequestService) SpecificValidation(c ResolverContext) error {
	return nil
}

func (dd *FullDIDDocRequestService) Prepare(c ResolverContext) error {
	return nil
}

func (dd *FullDIDDocRequestService) MakeAnswer(c ResolverContext) error {
	return nil
}
