package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
)

// ResourceDataEchoHandler godoc
//
//	@Summary		Fetch specific Resource
//	@Description	Get specific Resource within a DID Resource Collection
//	@Tags			Resource Resolution
//	@Accept			*/*
//	@Produce		*/*
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			resourceId	path		string	true	"Resource-specific unique-identifier"
//	@Success		200			{object}	[]byte
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Failure		501			{object}	types.IdentityError
//	@Router			/{did}/resources/{resourceId} [get]
func ResourceDataEchoHandler(c echo.Context) error {
	// Get Accept header
	contentType, profile := services.GetPriorityContentType(c.Request().Header.Get(echo.HeaderAccept), true)
	if contentType == types.JSONLD && profile == types.W3IDDIDURL {
		return services.EchoWrapHandler(&ResourceDataWithMetadataDereferencingService{})(c)
	}

	return services.EchoWrapHandler(&ResourceDataDereferencingService{})(c)
}

// ResourceMetadataEchoHandler godoc
//
//	@Summary		Fetch Resource-specific metadata
//	@Description	Get metadata for a specific Resource within a DID Resource Collection
//	@Tags			Resource Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			resourceId	path		string	true	"Resource-specific unique identifier"
//	@Success		200			{object}	types.DidDereferencing
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Failure		501			{object}	types.IdentityError
//	@Router			/{did}/resources/{resourceId}/metadata [get]
func ResourceMetadataEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&ResourceMetadataDereferencingService{})(c)
}
