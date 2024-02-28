package testconstants

import (
	"crypto/sha256"
	"fmt"
	"os"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func generateVerificationMethod() didTypes.VerificationMethod {
	return didTypes.VerificationMethod{
		Id:                     ExistentDid + "#key-1",
		VerificationMethodType: "JsonWebKey2020",
		Controller:             ExistentDid,
		VerificationMaterial:   ValidPubKeyJWK,
	}
}

func generateService() didTypes.Service {
	return didTypes.Service{
		Id:              ExistentDid + "#" + ValidServiceId,
		ServiceType:     "DIDCommMessaging",
		ServiceEndpoint: []string{"http://example.com"},
	}
}

func generateDIDDoc() didTypes.DidDoc {
	service := generateService()
	verificationMethod := generateVerificationMethod()

	return didTypes.DidDoc{
		Id:                 ExistentDid,
		VerificationMethod: []*didTypes.VerificationMethod{&verificationMethod},
		Service:            []*didTypes.Service{&service},
	}
}

func generateResource() []resourceTypes.ResourceWithMetadata {
	data := []byte("{\"attr\":[\"name\",\"age\"]}")
	checksum := sha256.New().Sum(data)
	return []resourceTypes.ResourceWithMetadata{
		{
			Resource: &resourceTypes.Resource{
				Data: data,
			},
			Metadata: &resourceTypes.Metadata{
				CollectionId: ValidIdentifier,
				Id:           ExistentResourceId,
				Name:         "Existing Resource Name",
				ResourceType: "string",
				MediaType:    "application/json",
				Checksum:     fmt.Sprintf("%x", checksum),
			},
		},
	}
}

func generateMetadata() didTypes.Metadata {
	return didTypes.Metadata{
		VersionId:   ValidVersionId,
		Deactivated: false,
		Created:     timestamppb.New(ValidCreated),
	}
}

func generateChecksum(data []byte) string {
	h := sha256.New()
	h.Write(data)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func getTestHostAddress() string {
	host := os.Getenv("TEST_HOST_ADDRESS")
	if host != "" {
		return host
	} else {
		return "localhost:8080"
	}
}
