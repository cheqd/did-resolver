package diddoc

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
)

type BaseDidDocRequestService struct {
	did                  string
	version              string
	fragment             string
	queries              url.Values
	result               types.DereferencingResultI
	requestedContentType types.ContentType
}

func (dd BaseDidDocRequestService) BasicValidation(c services.ResolverContext) error {
	didMethod, _, _, _ := types.TrySplitDID(dd.did)
	if didMethod != types.DID_METHOD {
		return types.NewMethodNotSupportedError(dd.did, dd.requestedContentType, nil, false)
	}

	err := utils.ValidateDID(dd.did, "", c.LedgerService.GetNamespaces())
	if err != nil {
		return types.NewInvalidDIDError(dd.did, dd.requestedContentType, nil, false)
	}

	return nil
}

func (dd *BaseDidDocRequestService) BasicPrepare(c services.ResolverContext) error {
	// Get DID from request
	did, err := services.GetDidParam(c)
	if err != nil {
		return types.NewInvalidDIDUrlError(c.Param("did"), dd.requestedContentType, err, true)
	}

	did = strings.Split(did, "#")[0]
	dd.requestedContentType = services.GetContentType(c.Request().Header.Get(echo.HeaderAccept))
	did, err = url.QueryUnescape(did)
	if err != nil {
		return types.NewInvalidDIDUrlError(did, dd.requestedContentType, err, true)
	}
	dd.did = did

	// Get Version
	dd.version = c.Param("version")

	return nil
}

func (dd *BaseDidDocRequestService) IsRedirectNeeded(c services.ResolverContext) bool {
	if !utils.IsValidDID(dd.did, "", c.LedgerService.GetNamespaces()) {
		err := utils.ValidateDID(dd.did, "", c.LedgerService.GetNamespaces())
		_, _, identifier, _ := types.TrySplitDID(dd.did)
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			return true
		}
	}
	return false
}

func (dd BaseDidDocRequestService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.did)
	queryRaw, _ := services.PrepareQueries(c)

	path := types.RESOLVER_PATH + migratedDid + utils.GetQuery(queryRaw) + utils.GetFragment(dd.fragment)
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *BaseDidDocRequestService) Query(c services.ResolverContext) error {
	result, err := c.DidDocService.Resolve(dd.did, dd.version, dd.requestedContentType)
	if err != nil {
		err.IsDereferencing = false
		return err
	}
	dd.result = result
	return nil
}

func (dd BaseDidDocRequestService) Respond(c services.ResolverContext) error {
	c.Response().Header().Set(echo.HeaderContentType, dd.result.GetContentType())
	return c.JSONPretty(http.StatusOK, dd.result, "  ")
}
