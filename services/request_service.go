package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"

	"github.com/rs/zerolog/log"

	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"google.golang.org/protobuf/runtime/protoiface"
)

type RequestService struct {
	didMethod                  string
	ledgerService              LedgerServiceI
	didDocService              DIDDocService
	resourceDereferenceService ResourceDereferenceService
}

func NewRequestService(didMethod string, ledgerService LedgerServiceI) RequestService {
	didDocService := DIDDocService{}
	return RequestService{
		didMethod:                  didMethod,
		ledgerService:              ledgerService,
		didDocService:              didDocService,
		resourceDereferenceService: NewResourceDereferenceService(ledgerService, didDocService),
	}
}

func (rs RequestService) ProcessDIDRequest(didUrl string, resolutionOptions types.ResolutionOption) ([]byte, int, types.ContentType) {
	var result []byte
	var statusCode int
	contentType := resolutionOptions.Accept
	if utils.IsDidUrl(didUrl) {
		log.Trace().Msgf("Dereferencing %s", didUrl)
		result, statusCode, contentType = rs.prepareDereferencingResult(didUrl, types.DereferencingOption(resolutionOptions))
	} else {
		log.Trace().Msgf("Resolving %s", didUrl)
		result, statusCode = rs.prepareResolutionResult(didUrl, resolutionOptions)
	}
	return result, statusCode, contentType
}

func (rs RequestService) prepareResolutionResult(did string, resolutionOptions types.ResolutionOption) ([]byte, int) {
	didResolution := rs.Resolve(did, resolutionOptions)

	resolutionMetadata, mErr1 := json.Marshal(didResolution.ResolutionMetadata)
	didDoc, mErr2 := rs.didDocService.MarshallDID(didResolution.Did)
	metadata, mErr3 := json.Marshal(&didResolution.Metadata)
	if mErr1 != nil || mErr2 != nil || mErr3 != nil {
		log.Error().Errs("errors", []error{mErr1, mErr2, mErr3}).Msg("Marshalling error")
		return createJsonResolutionInternalError(resolutionMetadata)
	}

	if didResolution.ResolutionMetadata.ResolutionError != "" {
		didDoc, metadata = "", []byte{}
	}

	result, err := createJsonResolution(didDoc, string(metadata), string(resolutionMetadata))
	if err != nil {
		log.Error().Err(err).Msg("Marshalling error")
		return createJsonResolutionInternalError([]byte{})
	}
	return result, didResolution.ResolutionMetadata.ResolutionError.GetStatusCode()
}

func (rs RequestService) prepareDereferencingResult(didUrl string, dereferencingOptions types.DereferencingOption) ([]byte, int, types.ContentType) {
	log.Info().Msgf("Dereferencing %s", didUrl)
	contentType := dereferencingOptions.Accept

	didDereferencing, statusCode := rs.Dereference(didUrl, dereferencingOptions)

	dereferencingMetadata, mErr1 := json.Marshal(didDereferencing.DereferencingMetadata)
	metadata, mErr2 := json.Marshal(didDereferencing.Metadata)
	if mErr1 != nil || mErr2 != nil {
		log.Error().Errs("errors", []error{mErr1, mErr2}).Msg("Marshalling error")
		response, errorStatusCode := createJsonDereferencingInternalError([]byte{})
		return response, errorStatusCode, contentType
	}

	if didDereferencing.DereferencingMetadata.ResolutionError != "" {
		didDereferencing.ContentStream = nil
		metadata = []byte{}
	} else {
		contentType = didDereferencing.DereferencingMetadata.ContentType
		if contentType != dereferencingOptions.Accept {
			return didDereferencing.ContentStream, statusCode, contentType
		}
	}

	result, err := createJsonDereferencing(didDereferencing.ContentStream, string(metadata), string(dereferencingMetadata))
	if err != nil {
		log.Error().Err(err).Msg("Marshalling error")
		response, errorStatusCode := createJsonDereferencingInternalError(dereferencingMetadata)
		return response, errorStatusCode, dereferencingOptions.Accept
	}
	return result, statusCode, contentType
}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (rs RequestService) Resolve(did string, resolutionOptions types.ResolutionOption) types.DidResolution {
	if !resolutionOptions.Accept.IsSupported() {
		return types.DidResolution{ResolutionMetadata: types.NewResolutionMetadata(did, types.JSON, types.RepresentationNotSupportedError)}
	}
	didResolutionMetadata := types.NewResolutionMetadata(did, resolutionOptions.Accept, "")

	if didMethod, _, _, _ := cheqdUtils.TrySplitDID(did); didMethod != rs.didMethod {
		didResolutionMetadata.ResolutionError = types.MethodNotSupportedError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}
	}
	if !cheqdUtils.IsValidDID(did, "", rs.ledgerService.GetNamespaces()) {
		didResolutionMetadata.ResolutionError = types.InvalidDIDError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}

	}

	didDoc, metadata, isFound, err := rs.ledgerService.QueryDIDDoc(did)
	if err != nil {
		didResolutionMetadata.ResolutionError = types.InternalError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}
	}

	resolvedMetadata, errorType := rs.ResolveMetadata(did, metadata)
	if errorType != "" {
		didResolutionMetadata.ResolutionError = types.InternalError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}
	}

	if !isFound {
		didResolutionMetadata.ResolutionError = types.NotFoundError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}
	}

	if didResolutionMetadata.ContentType == types.DIDJSONLD || didResolutionMetadata.ContentType == types.JSONLD {
		didDoc.Context = append(didDoc.Context, types.DIDSchemaJSONLD)
	} else {
		didDoc.Context = []string{}
	}
	return types.DidResolution{Did: didDoc, Metadata: resolvedMetadata, ResolutionMetadata: didResolutionMetadata}
}

// https://w3c-ccg.github.io/did-resolution/#dereferencing
func (rs RequestService) Dereference(didUrl string, dereferenceOptions types.DereferencingOption) (types.DidDereferencing, int) {
	did, path, query, fragmentId, err := cheqdUtils.TrySplitDIDUrl(didUrl)
	log.Info().Msgf("did: %s, path: %s, query: %s, fragmentId: %s", did, path, query, fragmentId)

	if !dereferenceOptions.Accept.IsSupported() {
		dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, dereferencingMetadata.ResolutionError.GetStatusCode()
	}

	if err != nil || !cheqdUtils.IsValidDIDUrl(didUrl, "", []string{}) {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.InvalidDIDUrlError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, dereferencingMetadata.ResolutionError.GetStatusCode()
	}

	// TODO: implement
	if query != "" {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, dereferencingMetadata.ResolutionError.GetStatusCode()
	}

	var didDereferencing types.DidDereferencing
	if path != "" {
		didDereferencing = rs.dereferencePrimary(path, did, didUrl, dereferenceOptions)
	} else {
		didDereferencing = rs.dereferenceSecondary(did, fragmentId, didUrl, dereferenceOptions)
	}

	return didDereferencing, didDereferencing.DereferencingMetadata.ResolutionError.GetStatusCode()
}

func (rs RequestService) dereferencePrimary(path string, did string, didUrl string, dereferenceOptions types.DereferencingOption) types.DidDereferencing {
	// Only resource are available for primary dereferencing
	return rs.resourceDereferenceService.DereferenceResource(path, did, didUrl, dereferenceOptions)
}

func (rs RequestService) dereferenceSecondary(did string, fragmentId string, didUrl string, dereferenceOptions types.DereferencingOption) types.DidDereferencing {
	didResolution := rs.Resolve(did, types.ResolutionOption(dereferenceOptions))

	dereferencingMetadata := types.DereferencingMetadata(didResolution.ResolutionMetadata)
	if dereferencingMetadata.ResolutionError != "" {
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	metadata := didResolution.Metadata

	var protoContent protoiface.MessageV1
	if fragmentId != "" {
		protoContent = rs.didDocService.GetDIDFragment(fragmentId, didResolution.Did)
		metadata = types.TransformToFragmentMetadata(metadata)
	} else {
		protoContent = &didResolution.Did
	}

	if protoContent == nil {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.NotFoundError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	jsonFragment, err := rs.didDocService.MarshallContentStream(protoContent, dereferenceOptions.Accept)
	if err != nil {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.InternalError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}
	contentStream := json.RawMessage(jsonFragment)

	return types.DidDereferencing{ContentStream: contentStream, Metadata: metadata, DereferencingMetadata: dereferencingMetadata}
}

func (rs RequestService) ResolveMetadata(did string, metadata cheqdTypes.Metadata) (types.ResolutionDidDocMetadata, types.ErrorType) {
	if metadata.Resources == nil {
		return types.NewResolutionDidDocMetadata(did, metadata, []*resourceTypes.ResourceHeader{}), ""
	}
	resources, errorType := rs.ledgerService.QueryCollectionResources(did)
	if errorType != "" {
		return types.ResolutionDidDocMetadata{}, errorType
	}
	return types.NewResolutionDidDocMetadata(did, metadata, resources), ""
}

func createJsonResolution(didDoc string, metadata string, resolutionMetadata string) ([]byte, error) {
	if didDoc == "" {
		didDoc = "null"
	}

	if metadata == "" {
		metadata = "[]"
	}

	response := struct {
		DidResolutionMetadata json.RawMessage `json:"didResolutionMetadata"`
		DidDocument           json.RawMessage `json:"didDocument"`
		DidDocumentMetadata   json.RawMessage `json:"didDocumentMetadata"`
	}{
		DidResolutionMetadata: json.RawMessage(resolutionMetadata),
		DidDocument:           json.RawMessage(didDoc),
		DidDocumentMetadata:   json.RawMessage(metadata),
	}

	respJson, err := json.Marshal(&response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal response")
		return []byte{}, err
	}

	return respJson, nil
}

func createJsonDereferencing(contentStream json.RawMessage, metadata string, dereferencingMetadata string) ([]byte, error) {
	if contentStream == nil {
		contentStream = json.RawMessage("null")
	}

	if metadata == "" {
		metadata = "[]"
	}

	response := struct {
		ContentStream         json.RawMessage `json:"contentStream"`
		ContentMetadata       json.RawMessage `json:"contentMetadata"`
		DereferencingMetadata json.RawMessage `json:"dereferencingMetadata"`
	}{
		ContentStream:         contentStream,
		ContentMetadata:       json.RawMessage(metadata),
		DereferencingMetadata: json.RawMessage(dereferencingMetadata),
	}

	respJson, err := json.Marshal(&response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal response")
		return []byte{}, err
	}

	return respJson, nil
}

func createJsonDereferencingInternalError(dereferencingMetadata []byte) ([]byte, int) {
	result, mErr := createJsonDereferencing(nil, "", string(dereferencingMetadata))
	if mErr != nil {
		return []byte{}, types.InternalError.GetStatusCode()
	}
	return result, types.InternalError.GetStatusCode()
}

func createJsonResolutionInternalError(resolutionMetadata []byte) ([]byte, int) {
	result, mErr := createJsonResolution("", "", string(resolutionMetadata))
	if mErr != nil {
		return []byte{}, types.InternalError.GetStatusCode()
	}
	return result, types.InternalError.GetStatusCode()
}
