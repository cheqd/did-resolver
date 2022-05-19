package services

import (

	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"
	"strings"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/golang/protobuf/jsonpb" //nolint
	"google.golang.org/protobuf/runtime/protoiface"
)

type DIDDocService struct{}

const (
	verificationMethod = "verificationMethod"
	publicKeyJwk       = "publicKeyJwk"
	didContext         = "context"
)

func (DIDDocService) MarshallProto(protoObject protoiface.MessageV1) (string, error) {
	var m jsonpb.Marshaler
	jsonObject, err := m.MarshalToString(protoObject)
	if err != nil {
		return "", err
	}
	return jsonObject, nil
}

func (ds DIDDocService) MarshallDID(didDoc cheqd.Did) (string, error) {
	mapDID, err := ds.protoToMap(&didDoc)
	if err != nil {
		return "", err
	}

	// VerKey changes
	formatedVerificationMethod, err := ds.MarshallVerificationMethod(didDoc.VerificationMethod)
	if err != nil {
		return "", err
	}
	mapDID[verificationMethod] = formatedVerificationMethod

	// Context changes
	if val, ok := mapDID[didContext]; ok {
		mapDID["@"+didContext] = val
		delete(mapDID, didContext)
	}

	result, err := json.Marshal(mapDID)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (ds DIDDocService) MarshallContentStream(contentStream protoiface.MessageV1, contentType types.ContentType) (string, error) {
	var mapContent map[string]interface{}
	var err error

	// VerKey changes, marshal
	if verificationMethod, ok := contentStream.(*cheqd.VerificationMethod); ok {
		mapContent, err = ds.prepareJWKPubkey(verificationMethod)
	} else {
		mapContent, err = ds.protoToMap(contentStream)
	}
	if err != nil {
		return "", err
	}

	// Context changes
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		mapContent["@"+didContext] = types.DIDSchemaJSONLD
	}

	result, err := json.Marshal(mapContent)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (DIDDocService) GetDIDFragment(fragmentId string, didDoc cheqd.Did) protoiface.MessageV1 {
	for _, verMethod := range didDoc.VerificationMethod {
		if strings.Contains(verMethod.Id, fragmentId) {
			return verMethod
		}
	}
	for _, service := range didDoc.Service {
		if strings.Contains(service.Id, fragmentId) {
			return service
		}
	}

	return nil
}

func (ds DIDDocService) prepareJWKPubkey(verificationMethod *cheqd.VerificationMethod) (map[string]interface{}, error) {
	methodJson, err := ds.protoToMap(verificationMethod)
	if err != nil {
		return nil, err
	}
	if len(verificationMethod.PublicKeyJwk) > 0 {
		methodJson[publicKeyJwk] = cheqd.PubKeyJWKToMap(verificationMethod.PublicKeyJwk)
	}
	return methodJson, nil
}

func (ds DIDDocService) MarshallVerificationMethod(verificationMethod []*cheqd.VerificationMethod) ([]map[string]interface{}, error) {
	var verMethodList []map[string]interface{}
	for _, value := range verificationMethod {
		methodJson, err := ds.prepareJWKPubkey(value)
		if err != nil {
			return nil, err
		}
		verMethodList = append(verMethodList, methodJson)
	}
	return verMethodList, nil
}

func (ds DIDDocService) protoToMap(protoObject protoiface.MessageV1) (map[string]interface{}, error) {
	jsonObj, err := ds.MarshallProto(protoObject)
	if err != nil {
		return nil, err
	}
	var mapObj map[string]interface{}

	err = json.Unmarshal([]byte(jsonObj), &mapObj)
	if err != nil {
		return nil, err
	}

	return mapObj, err
}
