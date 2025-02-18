package diddoc

import (
	"errors"
	"net/http"
	"strings"

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
	lowerQuery := strings.ToLower(dd.ResourceQuery)
	if lowerQuery != "true" && lowerQuery != "false" && lowerQuery != "" {
		return errors.New("invalid value for ResourceQuery: must be 'true' or 'false'")
	}
	return nil
}

func (dd *OnlyDIDDocRequestService) SpecificPrepare(c services.ResolverContext) error {
	return nil
}

func (dd OnlyDIDDocRequestService) Respond(c services.ResolverContext) error {
	_result := dd.Result.(*types.DidResolution)
	if dd.ResourceQuery == "false" {
		_result.Metadata = types.ResolutionDidDocMetadata{}
	}
	return c.JSONPretty(http.StatusOK, _result, "  ")
}
