package services

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
)

type (
	BaseRequestService struct {
		Did                  string
		Version              string
		Fragment             string
		IsDereferencing      bool
		Queries              url.Values
		Result               types.ResolutionResultI
		RequestedContentType types.ContentType
		Profile              string
	}
)

// Getters
func (dd BaseRequestService) GetDid() string {
	return dd.Did
}

func (dd BaseRequestService) GetContentType() types.ContentType {
	return dd.RequestedContentType
}

func (dd BaseRequestService) GetQueryParam(name string) string {
	return dd.Queries.Get(name)
}

func (dd BaseRequestService) GetDereferencing() bool {
	return dd.IsDereferencing
}

// Basic implementation
func (dd *BaseRequestService) BasicPrepare(c ResolverContext) error {
	// isDereferencingOrFragment variable to decide if we need to check if the resource is dereferencing or fragment
	isDereferencingOrFragment := dd.IsDereferencing && dd.Fragment == ""
	// Get Accept header
	dd.RequestedContentType, dd.Profile = GetPriorityContentType(c.Request().Header.Get(echo.HeaderAccept), isDereferencingOrFragment)
	if !dd.GetContentType().IsSupported() {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), types.JSON, nil, dd.IsDereferencing)
	}

	// Get DID from request
	did, err := GetDidParam(c)
	if err != nil {
		return types.NewInvalidDidUrlError(c.Param("did"), dd.RequestedContentType, err, dd.IsDereferencing)
	}

	// Get Did
	did = strings.Split(did, "#")[0]
	dd.Did = did

	// Get queries (We need to check that queries are allowed only for /:did path)
	queryRaw, flag := PrepareQueries(c)
	queries, err := url.ParseQuery(queryRaw)
	if err != nil {
		return err
	}
	if flag != nil {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}
	dd.Queries = queries

	return nil
}

func (dd BaseRequestService) BasicValidation(c ResolverContext) error {
	didMethod, _, _, _ := utils.TrySplitDID(dd.GetDid())
	if didMethod != types.DID_METHOD {
		return types.NewMethodNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	err := utils.ValidateDID(dd.GetDid(), "", c.LedgerService.GetNamespaces())
	if err != nil {
		return types.NewInvalidDidError(dd.GetDid(), dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	return nil
}

func (dd *BaseRequestService) IsRedirectNeeded(c ResolverContext) bool {
	if !utils.IsValidDID(dd.GetDid(), "", c.LedgerService.GetNamespaces()) {
		err := utils.ValidateDID(dd.GetDid(), "", c.LedgerService.GetNamespaces())
		_, _, identifier, _ := utils.TrySplitDID(dd.GetDid())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			return true
		}
	}
	return false
}

func (dd BaseRequestService) Redirect(c ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.GetDid())
	queryRaw, _ := PrepareQueries(c)

	path := types.RESOLVER_PATH + migratedDid + utils.GetQuery(queryRaw) + utils.GetFragment(dd.Fragment)
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *BaseRequestService) Query(c ResolverContext) error {
	result, err := c.DidDocService.Resolve(dd.GetDid(), dd.Version, dd.GetContentType())
	if err != nil {
		err.IsDereferencing = dd.GetDereferencing()
		return err
	}
	return dd.SetResponse(result)
}

func (dd BaseRequestService) SetupResponse(c ResolverContext) error {
	responseHeader := dd.Result.GetContentType()
	if dd.Profile != "" && responseHeader == string(types.JSONLD) {
		responseHeader = dd.Result.GetContentType() + ";profile=" + dd.Profile
	}
	c.Response().Header().Set(echo.HeaderContentType, responseHeader)
	if utils.IsGzipAccepted(c) {
		c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")
	}
	return nil
}

func (dd BaseRequestService) Respond(c ResolverContext) error {
	result := dd.FormatMetadataContentType(c)
	return c.JSONPretty(http.StatusOK, result, "  ")
}

// FormatMetadataContentType sets the ContentType of the result based on the profile and content type
func (dd BaseRequestService) FormatMetadataContentType(c ResolverContext) types.ResolutionResultI {
	switch result := dd.Result.(type) {
	case *types.DidResolution:
		// Set ContentType to DIDRES if profile is W3IDDIDRES and ContentType is JSONLD
		if dd.GetContentType() == types.JSONLD && dd.Profile == types.W3IDDIDRES {
			result.ResolutionMetadata.ContentType = types.DIDRES
			return result
		}
	default:
		return result
	}
	return dd.Result
}

// Setters

// SetResponse sets the response result
func (dd *BaseRequestService) SetResponse(response types.ResolutionResultI) error {
	dd.Result = response
	return nil
}

// Helpers

// RespondWithResourceData responds with the resource data
func (dd *BaseRequestService) RespondWithResourceData(c ResolverContext) error {
	c.Response().Header().Set(echo.HeaderContentType, dd.Result.GetContentType())

	return c.Blob(http.StatusOK, dd.Result.GetContentType(), dd.Result.GetBytes())
}
