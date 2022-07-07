package services

import (

	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"
	"strings"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/golang/protobuf/jsonpb" //nolint
	"github.com/iancoleman/orderedmap"
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
	mapDID.Set(verificationMethod, json.RawMessage(formatedVerificationMethod))

	// Context changes
	if val, ok := mapDID.Get(didContext); ok {
		mapDID.Set("@"+didContext, val)
		mapDID.Delete(didContext)
		mapDID.Sort(func(a *orderedmap.Pair, b *orderedmap.Pair) bool {
			return a.Key() == "@"+didContext
		})
	}

	result, err := json.Marshal(mapDID)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (ds DIDDocService) MarshallContentStream(contentStream protoiface.MessageV1, contentType types.ContentType) (string, error) {
	var mapContent orderedmap.OrderedMap
	var err error
	var context types.ContentType
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.DIDSchemaJSONLD
	}

	// VerKey changes, marshal
	if verificationMethod, ok := contentStream.(*cheqd.VerificationMethod); ok {
		mapContent, err = ds.prepareJWKPubkey(verificationMethod)
	} else if didDoc, ok := contentStream.(*cheqd.Did); ok {
		didDoc.Context = []string{string(context)}
		jsonDid, err := ds.MarshallDID(*didDoc)
		if err != nil {
			return "", err
		}
		return string(jsonDid), nil
		// Resource changes, marshal
	} else if resource, ok := contentStream.(*resource.Resource); ok {
		dResource := types.DereferencedResource{
			Context:           []string{string(context)},
			CollectionId:      resource.Header.CollectionId,
			Id:                resource.Header.Id,
			Name:              resource.Header.Name,
			ResourceType:      resource.Header.ResourceType,
			MediaType:         resource.Header.MediaType,
			Created:           resource.Header.Created,
			Checksum:          resource.Header.Checksum,
			PreviousVersionId: resource.Header.PreviousVersionId,
			NextVersionId:     resource.Header.NextVersionId,
			Data:              resource.Data,
		}
		jsonResource, err := json.Marshal(dResource)

		if err != nil {
			return "", err
		}
		return string(jsonResource), nil

	} else {
		mapContent, err = ds.protoToMap(contentStream)
	}
	if err != nil {
		return "", err
	}

	// Context changes
	if context != "" {
		mapContent.Set("@"+didContext, context)
		mapContent.Sort(func(a *orderedmap.Pair, b *orderedmap.Pair) bool {
			return a.Key() == "@"+didContext
		})
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

func (ds DIDDocService) prepareJWKPubkey(verificationMethod *cheqd.VerificationMethod) (orderedmap.OrderedMap, error) {
	methodJson, err := ds.protoToMap(verificationMethod)
	if err != nil {
		return *orderedmap.New(), err
	}
	if len(verificationMethod.PublicKeyJwk) > 0 {
		jsonKey, err := cheqd.PubKeyJWKToJson(verificationMethod.PublicKeyJwk)
		if err != nil {
			return *orderedmap.New(), err
		}
		methodJson.Set(publicKeyJwk, json.RawMessage(jsonKey))
	}
	return methodJson, nil
}

func (ds DIDDocService) MarshallVerificationMethod(verificationMethod []*cheqd.VerificationMethod) ([]byte, error) {
	var verMethodList []orderedmap.OrderedMap
	for _, value := range verificationMethod {
		methodJson, err := ds.prepareJWKPubkey(value)
		if err != nil {
			return []byte{}, err
		}
		verMethodList = append(verMethodList, methodJson)
	}
	return json.Marshal(verMethodList)
}

func (ds DIDDocService) protoToMap(protoObject protoiface.MessageV1) (orderedmap.OrderedMap, error) {
	mapObj := orderedmap.New()
	jsonObj, err := ds.MarshallProto(protoObject)
	if err != nil {
		return *mapObj, err
	}

	//var mapObj map[string]interface{}

	err = json.Unmarshal([]byte(jsonObj), &mapObj)
	if err != nil {
		return *mapObj, err
	}

	return *mapObj, err
}
