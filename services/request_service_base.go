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
		Result               types.ResultI
		RequestedContentType types.ContentType
	}
)

func (dd *BaseRequestService) BasicPrepare(c ResolverContext) error {
	// Here we raise errors even they were caught while getting the data from context

	// Get Accept header
	dd.RequestedContentType = GetContentType(c.Request().Header.Get(echo.HeaderAccept))
	if !dd.RequestedContentType.IsSupported() {
		return types.NewRepresentationNotSupportedError(dd.Did, types.JSON, nil, dd.IsDereferencing)
	}

	// Get DID from request
	did, err := GetDidParam(c)
	if err != nil {
		return types.NewInvalidDIDUrlError(c.Param("did"), dd.RequestedContentType, err, dd.IsDereferencing)
	}

	// Get Did
	did = strings.Split(did, "#")[0]
	dd.Did = did

	return nil
}

func (dd BaseRequestService) BasicValidation(c ResolverContext) error {
	didMethod, _, _, _ := types.TrySplitDID(dd.Did)
	if didMethod != types.DID_METHOD {
		return types.NewMethodNotSupportedError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	err := utils.ValidateDID(dd.Did, "", c.LedgerService.GetNamespaces())
	if err != nil {
		return types.NewInvalidDIDError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing)
	}

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

func (dd BaseRequestService) SetupResponse(c ResolverContext) error {
	c.Response().Header().Set(echo.HeaderContentType, dd.Result.GetContentType())
	if utils.IsGzipAccepted(c) {
		c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")
	}
	return nil
}

func (dd BaseRequestService) Respond(c ResolverContext) error {
	return c.JSONPretty(http.StatusOK, dd.Result, "  ")
}
