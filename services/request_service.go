package services

type RequestService struct {
	ledgerService LedgerService
	didDocService DIDDocService
}

func (r RequestService) processDIDRequest(did string, params map[string]string) string {
	didResolution := r.resolve(did, params)

	return ""
}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (RequestService) resolve(did string, resolutionOptions map[string]string) (string, string, string) {
	didDoc, metadata, err := LedgerService.QueryDIDDoc()
	didResolutionMetadata := 
	if err != nil {
		didResolutionMetadata = ResolutionErr(ResolutionNotFound)
	}
	return DidResolution{didResolutionMetadata, didDoc, metadata}
}

// https://w3c-ccg.github.io/did-resolution/#dereferencing
func (RequestService) dereference(didUrl string, dereferenceOptions map[string]string) (string, string, string) {
	did, metadata, err := LedgerService.QueryDIDDoc()
	if err != nil {
		didResolutionMetadata = ResolutionErr(ResolutionNotFound)
	}
	return dereferencingMetadata, contentStream, contentMetadata
}
