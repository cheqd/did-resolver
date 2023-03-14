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
		Queries              url.Values
		Result               types.ResultI
		RequestedContentType types.ContentType
	}
)

func (dd BaseRequestService) BasicValidation(c ResolverContext) error {
	didMethod, _, _, _ := types.TrySplitDID(dd.Did)
	if didMethod != types.DID_METHOD {
		return types.NewMethodNotSupportedError(dd.Did, dd.RequestedContentType, nil, false)
	}

	err := utils.ValidateDID(dd.Did, "", c.LedgerService.GetNamespaces())
	if err != nil {
		return types.NewInvalidDIDError(dd.Did, dd.RequestedContentType, nil, false)
	}

	if !dd.RequestedContentType.IsSupported() {
		return types.NewRepresentationNotSupportedError(dd.Did, types.JSON, nil, true)
	}

	return nil
}

func (dd *BaseRequestService) BasicPrepare(c ResolverContext) error {
	// Get DID from request
	did, err := GetDidParam(c)
	if err != nil {
		return types.NewInvalidDIDUrlError(c.Param("did"), dd.RequestedContentType, err, true)
	}

	did = strings.Split(did, "#")[0]
	dd.RequestedContentType = GetContentType(c.Request().Header.Get(echo.HeaderAccept))
	did, err = url.QueryUnescape(did)
	if err != nil {
		return types.NewInvalidDIDUrlError(did, dd.RequestedContentType, err, true)
	}
	dd.Did = did

	// Get Version
	dd.Version = c.Param("version")

	return nil
}

func (dd *BaseRequestService) IsRedirectNeeded(c ResolverContext) bool {
	if !utils.IsValidDID(dd.Did, "", c.LedgerService.GetNamespaces()) {
		err := utils.ValidateDID(dd.Did, "", c.LedgerService.GetNamespaces())
		_, _, identifier, _ := types.TrySplitDID(dd.Did)
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			return true
		}
	}
	return false
}

func (dd BaseRequestService) Redirect(c ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.Did)
	queryRaw, _ := PrepareQueries(c)

	path := types.RESOLVER_PATH + migratedDid + utils.GetQuery(queryRaw) + utils.GetFragment(dd.Fragment)
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *BaseRequestService) Query(c ResolverContext) error {
	result, err := c.DidDocService.Resolve(dd.Did, dd.Version, dd.RequestedContentType)
	if err != nil {
		err.IsDereferencing = false
		return err
	}
	dd.Result = result
	return nil
}

func (dd BaseRequestService) Respond(c ResolverContext) error {
	c.Response().Header().Set(echo.HeaderContentType, dd.Result.GetContentType())
	return c.JSONPretty(http.StatusOK, dd.Result, "  ")
}
