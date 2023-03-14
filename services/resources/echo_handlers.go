package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/labstack/echo/v4"
)

// DereferenceResourceMetadata godoc
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
//	@Router			/{did}/resources/{resourceId}/metadata [get]
func ResourceDataEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&ResourceDataDereferencingService{})(c)
}

func ResourceMetadataEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&ResourceMetadataDereferencingService{})(c)
}

func ResourceCollectionEchoHandler(c echo.Context) error {
	return services.EchoWrapHandler(&ResourceCollectionDereferencingService{})(c)
}