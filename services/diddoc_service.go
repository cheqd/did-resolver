package services

import (

	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/golang/protobuf/jsonpb" //nolint
	"github.com/golang/protobuf/proto"
)

type DIDDocService struct {
}

const (
	verificationMethod = "verificationMethod"
	publicKeyJwk       = "publicKeyJwk"
	didContext         = "context"
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

	// VerKey changes
	formatedVerificationMethod, err := ds.prepareJWKPubkey(didDoc)
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

func (DIDDocService) GetDIDFragment(DIDDoc cheqd.Did) string {
	//TODO: implement for dereferencing
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
