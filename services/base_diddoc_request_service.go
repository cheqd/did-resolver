package services

import (
	// "net/http"
	"net/http"
	"net/url"
	"strings"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
	// "strings"
)

type BaseDidDocRequestService struct {
	did                  string
	version              string
	fragment             string
	queries              url.Values
	result               types.DereferencingResultI
	requestedContentType types.ContentType
}

func (dd BaseDidDocRequestService) BasicValidation(c ResolverContext) error {
	didMethod, _, _, _ := types.TrySplitDID(dd.did)
	if didMethod != types.DID_METHOD {
		return types.NewMethodNotSupportedError(dd.did, dd.requestedContentType, nil, false)
	}

	err := utils.ValidateDID(dd.did, "", c.DidDocService.ledgerService.GetNamespaces())
	if err != nil {
		return types.NewInvalidDIDError(dd.did, dd.requestedContentType, nil, false)
	}

	return nil
}

func (dd *BaseDidDocRequestService) BasicPrepare(c ResolverContext) error {
	// Get DID from request
	did := strings.Split(c.Param("did"), "#")[0]
	dd.requestedContentType = getContentType(c.Request().Header.Get(echo.HeaderAccept))
	did, err := url.QueryUnescape(did)
	if err != nil {
		return types.NewInvalidDIDUrlError(did, dd.requestedContentType, err, true)
	}
	dd.did = did

	// Get Version
	dd.version = c.Param("version")

	return nil
}

func (dd *BaseDidDocRequestService) IsRedirectNeeded(c ResolverContext) bool {
	if !utils.IsValidDID(dd.did, "", c.LedgerService.GetNamespaces()) {
		err := utils.ValidateDID(dd.did, "", c.LedgerService.GetNamespaces())
		_, _, identifier, _ := types.TrySplitDID(dd.did)
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			return true
		}
	}
	return false
}

func (dd BaseDidDocRequestService) Redirect(c ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.did)
	queryRaw, _ := prepareQueries(c)
	path := types.RESOLVER_PATH + migratedDid + utils.GetQuery(queryRaw) + utils.GetFragment(dd.fragment)
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd BaseDidDocRequestService) Respond(c ResolverContext) error {
	c.Response().Header().Set(echo.HeaderContentType, dd.result.GetContentType())
	return c.JSONPretty(http.StatusOK, dd.result, "  ")
}

func (dd *BaseDidDocRequestService) Query(c ResolverContext) error {
	result, err := c.DidDocService.Resolve(dd.did, dd.version, dd.requestedContentType)
	if err != nil {
		return err
	}
	dd.result = result
	return nil
}

func DidDocEchoHandler(c echo.Context) error {
	// ToDo: Make
	isFragment := len(strings.Split(c.Param("did"), "#")) > 1
	isQuery := len(c.Request().URL.Query()) > 0
	isFullDidDoc := !isQuery && !isFragment

	switch {
	case isFullDidDoc:
		return EchoWrapHandler(&FullDIDDocRequestService{})(c)
	case isFragment:
		return EchoWrapHandler(&FragmentDIDDocRequestService{})(c)
	case isQuery:
		return EchoWrapHandler(&QueryDIDDocRequestService{})(c)
	default:
		// ToDo: make it more clearly
		return echo.NewHTTPError(500, "Invalid request")
	}
}

// Copy-paste
func DidDocVersionEchoHandler(c echo.Context) error {
	// ToDo: Make
	isFragment := c.Request().URL.Fragment != ""
	isQuery := len(c.Request().URL.Query()) > 0
	isFullDidDoc := !isQuery && !isFragment

	switch {
	case isFullDidDoc:
		return EchoWrapHandler(&FullDIDDocRequestService{})(c)
	case isFragment:
		return EchoWrapHandler(&FragmentDIDDocRequestService{})(c)
	case isQuery:
		return EchoWrapHandler(&QueryDIDDocRequestService{})(c)
	default:
		// ToDo: make it more clearly
		return echo.NewHTTPError(500, "Invalid request")
	}
}
