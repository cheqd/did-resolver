//go:build unit

package unit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

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
	e.Router().Find("GET", strings.Split(request.RequestURI, "?")[0], context)
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

var MockLedger = NewMockLedgerService(
	&testconstants.ValidDIDDoc,
	[]*didTypes.Metadata{&testconstants.ValidMetadata},
	[]resourceTypes.ResourceWithMetadata(testconstants.ValidResource),
)

type MockLedgerService struct {
	Did       *didTypes.DidDoc
	Metadata  []*didTypes.Metadata
	Resources []resourceTypes.ResourceWithMetadata
}

func NewMockLedgerService(did *didTypes.DidDoc, metadata []*didTypes.Metadata, resources []resourceTypes.ResourceWithMetadata) MockLedgerService {
	return MockLedgerService{
		Did:       did,
		Metadata:  metadata,
		Resources: resources,
	}
}

// TODO: add more unit tests for testing QueryDIDDoc method.
func (ls MockLedgerService) QueryDIDDoc(did string, version string) (*didTypes.DidDocWithMetadata, *types.IdentityError) {
	if ls.Did.Id == did {
		if version == "" {
			return &didTypes.DidDocWithMetadata{DidDoc: ls.Did, Metadata: ls.Metadata[len(ls.Metadata) -1]}, nil
		}
		for _, metadata := range ls.Metadata {
			if metadata.VersionId == version {
				return &didTypes.DidDocWithMetadata{DidDoc: ls.Did, Metadata: metadata}, nil
			}
		}
	}

	return nil, types.NewNotFoundError(did, types.JSON, nil, true)
}

// TODO: add unit tests for testing QueryAllDidDocVersionsMetadata method.
func (ls MockLedgerService) QueryAllDidDocVersionsMetadata(did string) ([]*didTypes.Metadata, *types.IdentityError) {
	if ls.Did.Id == did {
		return ls.Metadata, nil
	}

	return nil, types.NewNotFoundError(did, types.JSON, nil, true)
}

func (ls MockLedgerService) QueryResource(did string, resourceId string) (*resourceTypes.ResourceWithMetadata, *types.IdentityError) {
	if ls.Did.Id != did {
		return nil, types.NewNotFoundError(did, types.JSON, nil, true)
	}

	for _, resource := range ls.Resources {
		if resource.Metadata == nil {
			return nil, types.NewNotFoundError(did, types.JSON, nil, true)
		}
		if resource.Metadata.Id == resourceId {
			return &resource, nil
		}
	}

	return &resourceTypes.ResourceWithMetadata{}, types.NewNotFoundError(did, types.JSON, nil, true)
}

// TODO: add unit tests for testing QueryCollectionResources method.
func (ls MockLedgerService) QueryCollectionResources(did string) ([]*resourceTypes.Metadata, *types.IdentityError) {
	if ls.Did.Id != did {
		return []*resourceTypes.Metadata{}, types.NewNotFoundError(did, types.JSON, nil, true)
	}

	var metadata_list = make([]*resourceTypes.Metadata, len(ls.Resources))
	for i, resource := range ls.Resources {
		metadata_list[i] = resource.Metadata
	}

	return metadata_list, nil
}

func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}

func MustParseDate(sdate string) time.Time {
	date, err := time.Parse(time.RFC3339, sdate)
	if err != nil {
		panic(err)
	}

	return date
}
