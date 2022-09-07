package cmd

import (
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}
	identityError := generateIdentityError(err)
	if identityError.Code == 500 {
		log.Error().Err(identityError.Internal)
	} else {
		log.Warn().Err(identityError.Internal)
	}
	c.Response().Header().Set(echo.HeaderContentType, string(identityError.ContentType))
	err = c.JSONPretty(identityError.Code, identityError.DispalayMessage(), "  ")
	if err != nil {
		log.Error().Err(err)
	}
}

func generateIdentityError(err error) *types.IdentityError {
	identityError, ok := err.(*types.IdentityError)
	if ok {
		return identityError
	}
	he, ok := err.(*echo.HTTPError)
	if !ok || he.Code != 404 {
		return types.NewInternalError("", types.JSON, err, false)
	}
	return types.NewInvalidDIDUrlError("", types.JSON, err, false)
}
