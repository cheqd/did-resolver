package diddoc

import (
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	// Routes
	// Did docs
	e.GET(types.RESOLVER_PATH+":did", DidDocEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_METADATA, DidDocMetadataEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_VERSION_PATH+":version", DidDocVersionEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_VERSION_PATH+":version/metadata", DidDocVersionMetadataEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_VERSIONS_PATH, DidDocAllVersionMetadataEchoHandler)
}
