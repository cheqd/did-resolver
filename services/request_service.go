package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"
	"fmt"

	"github.com/cheqd/cheqd-did-resolver/types"
)

type RequestService struct {
	ledgerService LedgerService
	didDocService DIDDocService
}

func NewRequestService() RequestService {
	return RequestService{
		ledgerService: LedgerService{},
		didDocService: DIDDocService{},
	}
}

func (rs RequestService) ProcessDIDRequest(did string, params map[string]string) string {
	didResolution := rs.resolve(did, types.ResolutionOption{Accept: params["Accept"]})
	resolutionMetadata, err1 := json.Marshal(didResolution.ResolutionMetadata)
	didDoc, err2 := rs.didDocService.MarshallDID(didResolution.Did)
	metadata, err3 := rs.didDocService.MarshallProto(&didResolution.Metadata)

	if err1 != nil || err2 != nil || err3 != nil {
		resolutionMetadata := types.NewResolutionMetadata(resolutionOptions.Accept,
			types.ResolutionRepresentationNotSupported)
		resolutionMetadataJson, _ := json.Marshal(didResolution.ResolutionMetadata)
		return createJsonResolution("null", "null", string(resolutionMetadataJson))
	}

	if didResolution.ResolutionMetadata.ResolutionError != nil {
		return createJsonResolution("null", "null", string(resolutionMetadata))
	}

	return createJsonResolution(didDoc, metadata, string(resolutionMetadata))

}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (rs RequestService) resolve(did string, resolutionOptions types.ResolutionOption) types.DidResolution {
	didDoc, metadata, err := rs.ledgerService.QueryDIDDoc(did)
	didResolutionMetadata := types.NewResolutionMetadata(resolutionOptions.Accept, "")
	result := types.DidResolution{didDoc, metadata, didResolutionMetadata}
	if err != nil {
		didResolutionMetadata.ResolutionError = types.ResolutionNotFound
		return types.DidResolution{nil, nil, didResolutionMetadata}
	}
	return types.DidResolution{didDoc, metadata, didResolutionMetadata}
}

// https://w3c-ccg.github.io/did-resolution/#dereferencing
// func (RequestService) dereference(didUrl string, dereferenceOptions map[string]string) (string, string, string) {
// 	did, metadata, err := LedgerService.QueryDIDDoc()
// 	if err != nil {
// 		didResolutionMetadata = ResolutionErr(ResolutionNotFound)
// 	}
// 	return dereferencingMetadata, contentStream, contentMetadata
// }

func createJsonResolution(didDoc string, metadata string, resolutionMetadata string) {
	return fmt.Sprintf("{\"didDocument\" : %s,\"didDocumentMetadata\" : %s,\"didResolutionMetadata\" : %s}",
		didDoc, metadata, resolutionMetadata), ""
}
