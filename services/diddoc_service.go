package services

import (

	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"

	"github.com/cheqd/cheqd-did-resolver/types"
	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/golang/protobuf/jsonpb" //nolint
	"github.com/golang/protobuf/proto"
)

type DIDDocService struct {
}

const (
	verificationMethod = "verificationMethod"
	publicKeyJwk       = "publicKeyJwk"
)

func (DIDDocService) MarshallProto(protoObject proto.Message) (string, error) {
	var m jsonpb.Marshaler
	jsonObject, err := m.MarshalToString(protoObject)
	if err != nil {
		return "", err
	}
	return jsonObject, nil
}

func (ds DIDDocService) MarshallDID(didDoc cheqd.Did) (string, error) {
	jsonDID, err := ds.MarshallProto(&didDoc)
	if err != nil {
		return "", err
	}
	var mapDID map[string]interface{}
	json.Unmarshal([]byte(jsonDID), &mapDID)

	formatedVerificationMethod, err := ds.prepareJWKPubkey(didDoc)
	if err != nil {
		return "", err
	}

	mapDID[verificationMethod] = formatedVerificationMethod

	result, err := json.Marshal(mapDID)
	if err != nil {
		return "", err
	}
	return string(result), nil
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

func (ds DIDDocService) prepareJWKPubkey(didDoc cheqd.Did) ([]map[string]interface{}, error) {
	verMethodList := []map[string]interface{}{}
	for _, value := range didDoc.GetVerificationMethod() {
		methodJson, err := ds.protoToMap(value)
		if err != nil {
			return nil, err
		}
		if len(value.PublicKeyJwk) > 0 {
			methodJson[publicKeyJwk] = cheqd.PubKeyJWKToMap(value.PublicKeyJwk)
		}
		verMethodList = append(verMethodList, methodJson)

	}
	return verMethodList, nil
}

func (ds DIDDocService) protoToMap(protoObject proto.Message) (map[string]interface{}, error) {
	jsonObj, err := ds.MarshallProto(protoObject)
	if err != nil {
		return nil, err
	}
	var mapObj map[string]interface{}
	json.Unmarshal([]byte(jsonObj), &mapObj)
	return mapObj, err
}
