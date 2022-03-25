package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"github.com/cheqd/cheqd-did-resolver/types"
	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/golang/protobuf/jsonpb" //nolint
	"github.com/golang/protobuf/proto"
)

type DIDDocService struct {
}

func (DIDDocService) Marshall(protoObject proto.Message) (string, error) {
	var m jsonpb.Marshaler
	JsonObject, err := m.MarshalToString(protoObject)
	if err != nil {
		return "", err
	}
	return JsonObject, nil
}

func (DIDDocService) PrepareJWKPubkey(protoObject proto.Message) (string, error) {
}

func (DIDDocService) GetResolutionDIDMetadata(contentType string, errorType string) types.ResolutionMetadata {
	return types.ResolutionMetadata{}
}

func (DIDDocService) GetDID(contentType string, errorType string) cheqd.Did {
	return cheqd.Did{}
}

func (DIDDocService) GetDIDMetadata(contentType string, errorType string) cheqd.Metadata {
	return cheqd.Metadata{}
}

func (DIDDocService) GetDIDFragment(DIDDoc cheqd.Did) string {
	return ""
}
