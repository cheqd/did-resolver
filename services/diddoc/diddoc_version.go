package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

type DIDDocVersionRequestService struct {
	services.BaseRequestService
}

func (dd *DIDDocVersionRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = true
	return nil
}

func (dd *DIDDocVersionRequestService) SpecificPrepare(c services.ResolverContext) error {
	// Get Version
	dd.Version = c.Param("version")
	return nil
}

func (dd DIDDocVersionRequestService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.Did)

	path := types.RESOLVER_PATH + migratedDid + types.DID_VERSION_PATH + dd.Version
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *DIDDocVersionRequestService) SpecificValidation(c services.ResolverContext) error {
	if !utils.IsValidUUID(dd.Version) {
		return types.NewInvalidDIDUrlError(dd.Version, dd.GetContentType(), nil, dd.IsDereferencing)
	}
	return nil
}
