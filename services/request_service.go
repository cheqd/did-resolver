package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"

	"github.com/rs/zerolog/log"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"google.golang.org/protobuf/runtime/protoiface"
)

type RequestService struct {
	didMethod     string
	ledgerService LedgerServiceI
	didDocService DIDDocService
}

func NewRequestService(didMethod string, ledgerService LedgerServiceI) RequestService {
	return RequestService{
		didMethod:     didMethod,
		ledgerService: ledgerService,
		didDocService: DIDDocService{},
	}
}

func (rs RequestService) IsDidUrl(didUrl string) bool {
	_, path, query, fragmentId, err := cheqdUtils.TrySplitDIDUrl(didUrl)
	return err == nil && (path != "" || query != "" || fragmentId != "")
}

func (rs RequestService) ProcessDIDRequest(didUrl string, resolutionOptions types.ResolutionOption) (string, error) {
	if rs.IsDidUrl(didUrl) {
		log.Trace().Msgf("Dereferencing %s", didUrl)
		return rs.prepareDereferencingResult(didUrl, types.DereferencingOption(resolutionOptions))
	} else {
		log.Trace().Msgf("Resolving %s", didUrl)
		return rs.prepareResolutionResult(didUrl, resolutionOptions)
	}
}

func (rs RequestService) prepareResolutionResult(did string, resolutionOptions types.ResolutionOption) (string, error) {
	didResolution, err := rs.Resolve(did, resolutionOptions)
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

	return createJsonResolution(didDoc, metadata, string(resolutionMetadata))
}

func (rs RequestService) prepareDereferencingResult(did string, dereferencingOptions types.DereferencingOption) (string, error) {
	log.Info().Msgf("Dereferencing %s", did)

	didDereferencing, err := rs.Dereference(did, dereferencingOptions)
	if err != nil {
		return "", err
	}

	resolutionMetadata, err := json.Marshal(didDereferencing.DereferencingMetadata)
	if err != nil {
		return "", err
	}

	if didDereferencing.DereferencingMetadata.ResolutionError != "" {
		return createJsonDereferencing("", "", string(resolutionMetadata))
	}

	contentStream, err := rs.didDocService.MarshallContentStream(didDereferencing.ContentStream, dereferencingOptions.Accept)
	if err != nil {
		return "", err
	}

	metadata, err := rs.didDocService.MarshallProto(&didDereferencing.Metadata)
	if err != nil {
		return "", err
	}

	return createJsonDereferencing(contentStream, metadata, string(resolutionMetadata))
}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (rs RequestService) Resolve(did string, resolutionOptions types.ResolutionOption) (types.DidResolution, error) {
	didResolutionMetadata := types.NewResolutionMetadata(did, resolutionOptions.Accept, "")

	if didMethod, _, _, _ := cheqdUtils.TrySplitDID(did); didMethod != rs.didMethod {
		didResolutionMetadata.ResolutionError = types.ResolutionMethodNotSupported
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}, nil
	}

	if !cheqdUtils.IsValidDID(did, "", rs.ledgerService.GetNamespaces()) {
		didResolutionMetadata.ResolutionError = types.ResolutionInvalidDID
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

	if didResolutionMetadata.ContentType == types.DIDJSONLD {
		didDoc.Context = append(didDoc.Context, types.DIDSchemaJSONLD)
	}
	return types.DidResolution{Did: didDoc, Metadata: metadata, ResolutionMetadata: didResolutionMetadata}, nil
}

// https://w3c-ccg.github.io/did-resolution/#dereferencing
func (rs RequestService) Dereference(didUrl string, dereferenceOptions types.DereferencingOption) (types.DidDereferencing, error) {
	did, path, query, fragmentId, err := cheqdUtils.TrySplitDIDUrl(didUrl)
	log.Info().Msgf("did: %s, path: %s, query: %s, fragmentId: %s", did, path, query, fragmentId)

	// TODO: implement
	if query != "" {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.DereferencingNotSupported)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	if err != nil || !cheqdUtils.IsValidDIDUrl(didUrl, "", []string{}) {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.DereferencingInvalidDIDUrl)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	didResolution, err := rs.Resolve(did, types.ResolutionOption(dereferenceOptions))
	if err != nil {
		return types.DidDereferencing{}, err
	}
	dereferencingMetadata := types.DereferencingMetadata(didResolution.ResolutionMetadata)
	if dereferencingMetadata.ResolutionError != "" {
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}
	var protoiface.MessageV1 contentStream
	var protoiface.MessageV1 contentMetadata
	if path != "" {
		resourceId := utils.GetResourceId(path)
		// Only `resource` path is supported
		if resourceId == "" {
			dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.DereferencingNotSupported)
			return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
		}
		contentStream = rs.didDocService.GetResource(collectionId, resourceId)

	} else if fragmentId != "" {
		contentStream = rs.didDocService.GetDIDFragment(fragmentId, didResolution.Did)
		contentMetadata = didResolution.Metadata
	}

	if contentStream == nil {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.DereferencingNotFound)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	return types.DidDereferencing{ContentStream: contentStream, Metadata: contentMetadata, DereferencingMetadata: dereferencingMetadata}, nil
}

func createJsonResolution(didDoc string, metadata string, resolutionMetadata string) (string, error) {
	if didDoc == "" {
		didDoc = "null"
	}

	if metadata == "" {
		metadata = "[]"
	}

	response := struct {
		DidDocument           json.RawMessage `json:"didDocument"`
		DidDocumentMetadata   json.RawMessage `json:"didDocumentMetadata"`
		DidResolutionMetadata json.RawMessage `json:"didResolutionMetadata"`
	}{
		DidDocument:           json.RawMessage(didDoc),
		DidDocumentMetadata:   json.RawMessage(metadata),
		DidResolutionMetadata: json.RawMessage(resolutionMetadata),
	}

	respJson, err := json.Marshal(&response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal response")
		return "", err
	}

	return string(respJson), nil
}

func createJsonDereferencing(contentStream string, metadata string, dereferencingMetadata string) (string, error) {
	if contentStream == "" {
		contentStream = "null"
	}

	if metadata == "" {
		metadata = "[]"
	}

	response := struct {
		ContentStream         json.RawMessage `json:"contentStream"`
		ContentMetadata       json.RawMessage `json:"contentMetadata"`
		DereferencingMetadata json.RawMessage `json:"dereferencingMetadata"`
	}{
		ContentStream:         json.RawMessage(contentStream),
		ContentMetadata:       json.RawMessage(metadata),
		DereferencingMetadata: json.RawMessage(dereferencingMetadata),
	}

	respJson, err := json.Marshal(&response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal response")
		return "", err
	}

	return string(respJson), nil
}
