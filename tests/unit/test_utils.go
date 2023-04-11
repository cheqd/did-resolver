//go:build unit

package unit

import (
	"net/http"
	"net/http/httptest"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	didDocServices "github.com/cheqd/did-resolver/services/diddoc"
	resourceServices "github.com/cheqd/did-resolver/services/resource"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
)

func DefineContentType(expectedContentType types.ContentType, resolutionType types.ContentType) types.ContentType {
	if expectedContentType == "" {
		return resolutionType
	}

	return expectedContentType
}

func SetupEmptyContext(request *http.Request, resolutionType types.ContentType, ledgerService services.LedgerServiceI) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	didDocServices.SetRoutes(e)
	resourceServices.SetRoutes(e)

	didService := services.NewDIDDocService(types.DID_METHOD, ledgerService)
	resourceService := services.NewResourceService(types.DID_METHOD, ledgerService)

	rec := httptest.NewRecorder()
	context := e.NewContext(request, rec)
	e.Router().Find("GET", request.RequestURI, context)
	rc := services.ResolverContext{
		Context:         context,
		LedgerService:   ledgerService,
		DidDocService:   didService,
		ResourceService: resourceService,
	}

	request.Header.Add("accept", string(resolutionType))
	return rc, rec
}

type RedirectDIDTestCase struct {
	DidURL                 string
	ResolutionType         types.ContentType
	ExpectedDidURLRedirect string
	ExpectedError          error
}

var MockLedger = NewMockLedgerService(&testconstants.ValidDIDDoc, &testconstants.ValidMetadata, &testconstants.ValidResource)

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

// TODO: add more unit tests for testing QueryDIDDoc method.
func (ls MockLedgerService) QueryDIDDoc(did string, version string) (*didTypes.DidDocWithMetadata, *types.IdentityError) {
	if ls.Did.Id == did {
		return &didTypes.DidDocWithMetadata{DidDoc: ls.Did, Metadata: ls.Metadata}, nil
	}

	return nil, types.NewNotFoundError(did, types.JSON, nil, true)
}

// TODO: add unit tests for testing QueryAllDidDocVersionsMetadata method.
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

// TODO: add unit tests for testing QueryCollectionResources method.
func (ls MockLedgerService) QueryCollectionResources(did string) ([]*resourceTypes.Metadata, *types.IdentityError) {
	if ls.Did.Id != did || ls.Resource.Metadata == nil {
		return []*resourceTypes.Metadata{}, types.NewNotFoundError(did, types.JSON, nil, true)
	}

	return []*resourceTypes.Metadata{ls.Resource.Metadata}, nil
}

func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}
