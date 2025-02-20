package diddoc

import (
	"net/http"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type OnlyDIDDocRequestService struct {
	services.BaseRequestService
	ResourceQuery string
}

func (dd *OnlyDIDDocRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = false
	return nil
}

func (dd *OnlyDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	return nil
}

func (dd *OnlyDIDDocRequestService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}

func (dd OnlyDIDDocRequestService) Respond(c services.ResolverContext) error {
	_result := dd.Result.(*types.DidResolution)
	_result.Metadata = nil
	return c.JSONPretty(http.StatusOK, _result, "  ")
}
