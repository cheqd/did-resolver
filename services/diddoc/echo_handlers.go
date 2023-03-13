package diddoc

import (
	"strings"

	"github.com/cheqd/did-resolver/services"
	"github.com/labstack/echo/v4"
)

func DidDocEchoHandler(c echo.Context) error {
	// ToDo: Make fragment detection better
	isFragment := len(strings.Split(c.Param("did"), "#")) > 1
	isQuery := len(c.Request().URL.Query()) > 0
	isFullDidDoc := !isQuery && !isFragment

	switch {
	case isFullDidDoc:
		return services.EchoWrapHandler(&FullDIDDocRequestService{})(c)
	case isFragment:
		return services.EchoWrapHandler(&FragmentDIDDocRequestService{})(c)
	case isQuery:
		return services.EchoWrapHandler(&QueryDIDDocRequestService{})(c)
	default:
		// ToDo: make it more clearly
		return echo.NewHTTPError(500, "Unknown handler for processing request")
	}
}

func DidDocVersionEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&DIDDocVersionRequestService{})(c)
}

// ResolveDIDDocVersionMetadata godoc
//
//	@Summary		Resolve DID Document Version Metadata on did:cheqd
//	@Description	Fetch metadata of specific a DID Document ("DIDDoc") version for a given DID and version ID
//	@Tags			DID Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+jsonww
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			versionId	path		string	true	"version of a DID document"
//	@Success		200			{object}	types.DidDereferencing
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Router			/{did}/version/{versionId}/metadata [get]
func DidDocVersionMetadataEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&DIDDocVersionMetadataRequestService{})(c)
}

// ResolveAllDidDocVersionsMetadata godoc
//
//	@Summary		Resolve DID Document Version on did:cheqd
//	@Description	Fetch specific all version of a DID Document ("DIDDoc") for a given DID and version ID
//	@Tags			DID Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			versionId	path		string	true	"version of a DID document"
//	@Success		200			{object}	types.DidResolution
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Router			/{did}/version/{versionId} [get]
func DidDocAllVersionMetadataEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&DIDDocAllVersionMetadataRequestService{})(c)
}