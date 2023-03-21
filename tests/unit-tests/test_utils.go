package tests

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var mockLedgerService = NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource)

const (
	ValidMethod     = "cheqd"
	ValidNamespace  = "mainnet"
	ValidIdentifier = "fb53dd05-329b-4614-a3f2-c0a8c7554ee3"
	ValidDid        = "did:" + ValidMethod + ":" + ValidNamespace + ":" + ValidIdentifier
	ValidResourceId = "a09abea0-22e0-4b35-8f70-9cc3a6d0b5fd"
	ValidPubKeyJWK  = "{" +
		"\"crv\":\"Ed25519\"," +
		"\"kid\":\"_Qq0UL2Fq651Q0Fjd6TvnYE-faHiOpRlPVQcY_-tA4A\"," +
		"\"kty\":\"OKP\"," +
		"\"x\":\"VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ\"" +
		"}"
	ValidVersionId = "test_version_id"
)

const (
	InvalidMethod     = "invalid_method"
	InvalidNamespace  = "invalid_namespace"
	InvalidIdentifier = "invalid_identifier"
	InvalidDid        = "did:" + InvalidMethod + ":" + InvalidNamespace + ":" + InvalidIdentifier
	InvalidResourceId = "invalid_resource_id"
)

const (
	NotExistIdentifier = "fb53dd05-329b-4614-a3f2-c0a8c7ffffff"
	NotExistDID        = "did:" + ValidMethod + ":" + ValidNamespace + ":" + NotExistIdentifier
)

var (
	EmptyTimestamp = &timestamppb.Timestamp{
		Seconds: 0,
		Nanos:   0,
	}
	EmptyTime = EmptyTimestamp.AsTime()

	NotEmptyTimestamp = &timestamppb.Timestamp{
		Seconds: 123456789,
		Nanos:   0,
	}
	NotEmptyTime = NotEmptyTimestamp.AsTime()
)

var (
	ResourceData     = []byte("test_checksum")
	ResourceMetadata = resourceTypes.Metadata{
		CollectionId: ValidIdentifier,
		Id:           ValidResourceId,
		Name:         "Existing Resource Name",
		ResourceType: "CL-Schema",
		MediaType:    "application/json",
		Checksum:     generateChecksum(ResourceData),
	}

	ValidMetadataResource = types.DereferencedResource{
		ResourceURI:       ValidDid + types.RESOURCE_PATH + ResourceMetadata.Id,
		CollectionId:      ResourceMetadata.CollectionId,
		ResourceId:        ResourceMetadata.Id,
		Name:              ResourceMetadata.Name,
		ResourceType:      ResourceMetadata.ResourceType,
		MediaType:         ResourceMetadata.MediaType,
		Created:           &EmptyTime,
		Checksum:          ResourceMetadata.Checksum,
		PreviousVersionId: nil,
		NextVersionId:     nil,
	}
)

var (
	validDIDDoc             = ValidDIDDoc()
	validVerificationMethod = ValidVerificationMethod()
	validDIDDocResolution   = types.NewDidDoc(&validDIDDoc)
	validMetadata           = ValidMetadata()
	validService            = ValidService()
	validResource           = ValidResource()
	validFragmentMetadata   = types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{})
	validQuery, _           = url.ParseQuery("attr=value")
)

var (
	dereferencedResourceList = types.NewDereferencedResourceList(ValidDid, []*resourceTypes.Metadata{validResource.Metadata})
	resolutionDIDDocMetadata = types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{validResource.Metadata})
	resourceData             = types.DereferencedResourceData(validResource.Resource.Data)
)

func ValidVerificationMethod() didTypes.VerificationMethod {
	return didTypes.VerificationMethod{
		Id:                     ValidDid + "#key-1",
		VerificationMethodType: "JsonWebKey2020",
		Controller:             ValidDid,
		VerificationMaterial:   ValidPubKeyJWK,
	}
}

func ValidService() didTypes.Service {
	return didTypes.Service{
		Id:              ValidDid + "#service-1",
		ServiceType:     "DIDCommMessaging",
		ServiceEndpoint: []string{"http://example.com"},
	}
}

func ValidDIDDoc() didTypes.DidDoc {
	service := ValidService()
	verificationMethod := ValidVerificationMethod()

	return didTypes.DidDoc{
		Id:                 ValidDid,
		VerificationMethod: []*didTypes.VerificationMethod{&verificationMethod},
		Service:            []*didTypes.Service{&service},
	}
}

func ValidResource() resourceTypes.ResourceWithMetadata {
	data := []byte("{\"attr\":[\"name\",\"age\"]}")
	checksum := sha256.New().Sum(data)
	return resourceTypes.ResourceWithMetadata{
		Resource: &resourceTypes.Resource{
			Data: data,
		},
		Metadata: &resourceTypes.Metadata{
			CollectionId: ValidIdentifier,
			Id:           ValidResourceId,
			Name:         ValidResourceId,
			ResourceType: "string",
			MediaType:    "application/json",
			Checksum:     fmt.Sprintf("%x", checksum),
		},
	}
}

func ValidMetadata() didTypes.Metadata {
	return didTypes.Metadata{VersionId: "test_version_id", Deactivated: false}
}

func generateChecksum(data []byte) string {
	h := sha256.New()
	h.Write(data)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func defineContentType(expectedContentType types.ContentType, resolutionType types.ContentType) types.ContentType {
	if expectedContentType == "" {
		return resolutionType
	}

	return expectedContentType
}

func getDID(didURL string) string {
	return strings.Split(didURL, "/")[3]
}

func getResourceId(didURL string) string {
	return strings.Split(didURL, "/")[5]
}

func setupContext(path string, paramsNames []string, paramsValues []string, resolutionType types.ContentType, ledgerService services.LedgerServiceI) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()

	didService := services.NewDIDDocService(types.DID_METHOD, ledgerService)
	resourceService := services.NewResourceService(types.DID_METHOD, ledgerService)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	rc := services.ResolverContext{
		Context:         context,
		LedgerService:   ledgerService,
		DidDocService:   didService,
		ResourceService: resourceService,
	}
	rc.SetPath(path)
	rc.SetParamNames(paramsNames...)
	rc.SetParamValues(paramsValues...)
	req.Header.Add("accept", string(resolutionType))
	return rc, rec
}

type MockLedgerService struct {
	Did      *didTypes.DidDoc
	Metadata *didTypes.Metadata
	Resource *resourceTypes.ResourceWithMetadata
}

func NewMockLedgerService(did *didTypes.DidDoc, metadata *didTypes.Metadata, resource *resourceTypes.ResourceWithMetadata) MockLedgerService {
	return MockLedgerService{
		Did:      did,
		Metadata: metadata,
		Resource: resource,
	}
}

func (ls MockLedgerService) QueryDIDDoc(did string, version string) (*didTypes.DidDocWithMetadata, *types.IdentityError) {
	if ls.Did.Id == did {
		println("query !!!" + ls.Did.Id)
		return &didTypes.DidDocWithMetadata{DidDoc: ls.Did, Metadata: ls.Metadata}, nil
	}
	return nil, types.NewNotFoundError(did, types.JSON, nil, true)
}

func (ls MockLedgerService) QueryAllDidDocVersionsMetadata(did string) ([]*didTypes.Metadata, *types.IdentityError) {
	if ls.Did.Id == did {
		return []*didTypes.Metadata{ls.Metadata}, nil
	}

	return nil, types.NewNotFoundError(did, types.JSON, nil, true)
}

func (ls MockLedgerService) QueryResource(did string, resourceId string) (*resourceTypes.ResourceWithMetadata, *types.IdentityError) {
	if ls.Did.Id != did || ls.Resource.Metadata == nil || ls.Resource.Metadata.Id != resourceId {
		return nil, types.NewNotFoundError(did, types.JSON, nil, true)
	}
	return ls.Resource, nil
}

func (ls MockLedgerService) QueryCollectionResources(did string) ([]*resourceTypes.Metadata, *types.IdentityError) {
	if ls.Did.Id != did || ls.Resource.Metadata == nil {
		return []*resourceTypes.Metadata{}, types.NewNotFoundError(did, types.JSON, nil, true)
	}
	return []*resourceTypes.Metadata{ls.Resource.Metadata}, nil
}

func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}
