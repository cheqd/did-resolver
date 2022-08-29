package services

import (
	"strings"

	"github.com/cheqd/did-resolver/types"
)

type DIDDocService struct{}

func IsFragmentId(id string, requestedId string) bool {
	if strings.Contains(id, "#") {
		id = strings.Split(id, "#")[1]
	}
	return id == requestedId
}

func (DIDDocService) GetDIDFragment(fragmentId string, didDoc types.DidDoc) types.ContentStreamI {
	for _, verMethod := range didDoc.VerificationMethod {
		if IsFragmentId(verMethod.Id, fragmentId) {
			return &verMethod
		}
	}
	for _, service := range didDoc.Service {
		if IsFragmentId(service.Id, fragmentId) {
			return &service
		}
	}

	return nil
}

func (DIDDocService) GetDIDService(queryId string, didDoc cheqd.Did) *cheqd.Service {
	for _, service := range didDoc.Service {
		if IsFragmentId(service.Id, queryId) {
			return service
		}
	}
	return nil
}

func CreateServiceEndpoint(relativeRef string, fragmentId string, inputServiceEndpoint string) (outputServiceEndpoint string) {
	outputServiceEndpoint = inputServiceEndpoint + relativeRef
	if fragmentId != "" {
		outputServiceEndpoint += "#" + fragmentId
	}
	return outputServiceEndpoint
}
