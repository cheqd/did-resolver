package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/migrations"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DIDDocVersionRequestService struct {
	BaseDidDocRequestService
}

func (dd *DIDDocVersionRequestService) Prepare(c services.ResolverContext) error {
	return nil
}

func (dd DIDDocVersionRequestService) Redirect(c services.ResolverContext) error {
	migratedDid := migrations.MigrateDID(dd.did)

	path := types.RESOLVER_PATH + migratedDid + types.DID_VERSION_PATH + dd.version
	return c.Redirect(http.StatusMovedPermanently, path)
}

func (dd *DIDDocVersionRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *DIDDocVersionRequestService) MakeResponse(c services.ResolverContext) error {
	return nil
}
