package diddoc

import (
	"errors"
	"strings"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
)

// DidDocEchoHandler godoc
//
//	@Summary		Resolve DID Document on did:cheqd
//	@Description	Fetch DID Document ("DIDDoc") from cheqd network
//	@Tags			DID Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			service		query		string	false	"Service Type"
//	@Param			fragmentId	query		string	false	"#Fragment"
//	@Param			versionId	query		string	false	"Version"
//	@Success		200			{object}	types.DidResolution
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Failure		501			{object}	types.IdentityError
//	@Router			/{did} [get]
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
		return types.NewInternalError(c.Param("did"), types.JSON, errors.New("Unknown internal error while getting the type of query"), true)
	}
}

// DidDocVersionEchoHandler godoc
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
//	@Failure		501			{object}	types.IdentityError
//	@Router			/{did}/version/{versionId} [get]
func DidDocVersionEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&DIDDocVersionRequestService{})(c)
}

// DidDocVersionMetadataEchoHandler godoc
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
//	@Failure		501			{object}	types.IdentityError
//	@Router			/{did}/version/{versionId}/metadata [get]
func DidDocVersionMetadataEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&DIDDocVersionMetadataRequestService{})(c)
}

// DidDocAllVersionMetadataEchoHandler godoc
//
//	@Summary		Resolve DID Document Versions on did:cheqd
//	@Description	Fetch specific all versions of a DID Document ("DIDDoc") for a given DID
//	@Tags			DID Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did	path		string	true	"Full DID with unique identifier"
//	@Success		200	{object}	types.ResourceDereferencing{contentStream=types.DereferencedDidVersionsList}
//	@Failure		400	{object}	types.IdentityError
//	@Failure		404	{object}	types.IdentityError
//	@Failure		406	{object}	types.IdentityError
//	@Failure		500	{object}	types.IdentityError
//	@Failure		501	{object}	types.IdentityError
//	@Router			/{did}/versions [get]
func DidDocAllVersionMetadataEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&DIDDocAllVersionMetadataRequestService{})(c)
}
