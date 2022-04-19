package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"
	"fmt"

	"github.com/cheqd/cheqd-did-resolver/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/spf13/viper"
)

type RequestService struct {
	ledgerService LedgerServiceI
	didDocService DIDDocService
}

func NewRequestService(ledgerService LedgerServiceI) RequestService {
	return RequestService{
		ledgerService: ledgerService,
		didDocService: DIDDocService{},
	}
}

func (rs RequestService) ProcessDIDRequest(did string, resolutionOptions types.ResolutionOption) (string, error) {

	didResolution, err := rs.Resolve(did, types.ResolutionOption{resolutionOptions.Accept})
	if err != nil {
		return "", err
	}

	resolutionMetadata, err := json.Marshal(didResolution.ResolutionMetadata)
	if err != nil {
		return "", err
	}

	didDoc, err := rs.didDocService.MarshallDID(didResolution.Did)
	if err != nil {
		return "", err
	}

	metadata, err := rs.didDocService.MarshallProto(&didResolution.Metadata)
	if err != nil {
		return "", err
	}

	if didResolution.ResolutionMetadata.ResolutionError != "" {
		didDoc, metadata = "", ""
	}

	return createJsonResolution(didDoc, metadata, string(resolutionMetadata)), nil

}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (rs RequestService) Resolve(did string, resolutionOptions types.ResolutionOption) (types.DidResolution, error) {

	didResolutionMetadata := types.NewResolutionMetadata(did, resolutionOptions.Accept, "")

	method := viper.GetString("method")
	if !cheqdUtils.IsValidDID(did, method, rs.ledgerService.GetNamespaces()) {
		if didMethod, _, _, _ := cheqdUtils.TrySplitDID(did); didMethod != method {
			didResolutionMetadata.ResolutionError = types.ResolutionMethodNotSupported
		} else {
			didResolutionMetadata.ResolutionError = types.ResolutionInvalidDID
		}
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}, nil
	}

	didDoc, metadata, isFound, err := rs.ledgerService.QueryDIDDoc(did)
	if err != nil {
		return types.DidResolution{}, err
	}
	if !isFound {
		didResolutionMetadata.ResolutionError = types.ResolutionNotFound
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}, nil
	}
	if didResolutionMetadata.ContentType == types.ResolutionDIDJSONLDType {
		didDoc.Context = append(didDoc.Context, types.DIDSchemaJSONLD)
	}
	return types.DidResolution{didDoc, metadata, didResolutionMetadata}, nil
}

// https://w3c-ccg.github.io/did-resolution/#dereferencing
// func (RequestService) dereference(didUrl string, dereferenceOptions map[string]string) (string, string, string) {
// 	did, metadata, err := LedgerService.QueryDIDDoc()
// 	if err != nil {
// 		didResolutionMetadata = ResolutionErr(ResolutionNotFound)
// 	}
// 	return dereferencingMetadata, contentStream, contentMetadata
// }

func createJsonResolution(didDoc string, metadata string, resolutionMetadata string) string {
	if didDoc == "" {
		didDoc = "null"
	}
	if metadata == "" {
		metadata = "[]"
	}
	return fmt.Sprintf("{\"didDocument\" : %s,\"didDocumentMetadata\" : %s,\"didResolutionMetadata\" : %s}",
		didDoc, metadata, resolutionMetadata)
}
