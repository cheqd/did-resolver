package services

import (
	"strings"

	"github.com/cheqd/did-resolver/types"
)

type DIDDocService struct{}

func (DIDDocService) GetDIDFragment(fragmentId string, didDoc types.DidDoc) types.ContentStreamI {
	for _, verMethod := range didDoc.VerificationMethod {
		if strings.Contains(verMethod.Id, fragmentId) {
			return &verMethod
		}
	}
	for _, service := range didDoc.Service {
		if strings.Contains(service.Id, fragmentId) {
			return &service
		}
	}

	return nil
}
