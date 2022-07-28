package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"

	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
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
	var result string
	var err error
	if rs.IsDidUrl(didUrl) {
		log.Trace().Msgf("Dereferencing %s", didUrl)
		result, err = rs.prepareDereferencingResult(didUrl, types.DereferencingOption(resolutionOptions))
	} else {
		log.Trace().Msgf("Resolving %s", didUrl)
		result, err = rs.prepareResolutionResult(didUrl, resolutionOptions)
	}
	return result, err
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

	metadata, err := json.Marshal(&didResolution.Metadata)
	if err != nil {
		return "", err
	}

	if didResolution.ResolutionMetadata.ResolutionError != "" {
		didDoc, metadata = "", []byte{}
	}

	return createJsonResolution(didDoc, string(metadata), string(resolutionMetadata))
}

func (rs RequestService) prepareDereferencingResult(did string, dereferencingOptions types.DereferencingOption) (string, error) {
	log.Info().Msgf("Dereferencing %s", did)

	didDereferencing, err := rs.Dereference(did, dereferencingOptions)
	if err != nil {
		return "", err
	}

	dereferencingMetadata, err := json.Marshal(didDereferencing.DereferencingMetadata)
	if err != nil {
		return "", err
	}

	if didDereferencing.DereferencingMetadata.ResolutionError != "" {
		return createJsonDereferencing(nil, "", string(dereferencingMetadata))
	}

	metadata, err := json.Marshal(didDereferencing.Metadata)
	if err != nil {
		return "", err
	}

	return createJsonDereferencing(didDereferencing.ContentStream, string(metadata), string(dereferencingMetadata))
}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (rs RequestService) Resolve(did string, resolutionOptions types.ResolutionOption) (types.DidResolution, error) {
	if !resolutionOptions.Accept.IsSupported() {
		return types.DidResolution{ResolutionMetadata: types.NewResolutionMetadata(did, types.JSON, types.RepresentationNotSupportedError)}, nil
	}
	didResolutionMetadata := types.NewResolutionMetadata(did, resolutionOptions.Accept, "")

	if didMethod, _, _, _ := cheqdUtils.TrySplitDID(did); didMethod != rs.didMethod {
		didResolutionMetadata.ResolutionError = types.MethodNotSupportedError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}, nil
	}

	if !cheqdUtils.IsValidDID(did, "", rs.ledgerService.GetNamespaces()) {
		didResolutionMetadata.ResolutionError = types.InvalidDIDError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}, nil

	}

	didDoc, metadata, isFound, err := rs.ledgerService.QueryDIDDoc(did)
	if err != nil {
		return types.DidResolution{}, err
	}

	resolvedMetadata, err := rs.ResolveMetadata(did, metadata)
	if err != nil {
		return types.DidResolution{}, err
	}

	if !isFound {
		didResolutionMetadata.ResolutionError = types.NotFoundError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}, nil
	}

	if didResolutionMetadata.ContentType == types.DIDJSONLD || didResolutionMetadata.ContentType == types.JSONLD {
		didDoc.Context = append(didDoc.Context, types.DIDSchemaJSONLD)
	} else {
		didDoc.Context = []string{}
	}
	return types.DidResolution{Did: didDoc, Metadata: resolvedMetadata, ResolutionMetadata: didResolutionMetadata}, nil
}

// https://w3c-ccg.github.io/did-resolution/#dereferencing
func (rs RequestService) Dereference(didUrl string, dereferenceOptions types.DereferencingOption) (types.DidDereferencing, error) {
	did, path, query, fragmentId, err := cheqdUtils.TrySplitDIDUrl(didUrl)
	log.Info().Msgf("did: %s, path: %s, query: %s, fragmentId: %s", did, path, query, fragmentId)

	if err != nil || !cheqdUtils.IsValidDIDUrl(didUrl, "", []string{}) {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.InvalidDIDUrlError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	// TODO: implement
	if query != "" {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	var didDereferencing types.DidDereferencing
	if path != "" {
		didDereferencing, err = rs.dereferencePrimary(path, did, didUrl, dereferenceOptions)
	} else {
		didDereferencing, err = rs.dereferenceSecondary(did, fragmentId, didUrl, dereferenceOptions)
	}

	if err != nil {
		return types.DidDereferencing{}, err
	}

	return didDereferencing, nil
}

func (rs RequestService) dereferencePrimary(path string, did string, didUrl string, dereferenceOptions types.DereferencingOption) (types.DidDereferencing, error) {
	resourceId := utils.GetResourceId(path)
	// Only `resource` path is supported
	if resourceId == "" {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	resource, isFound, err := rs.ledgerService.QueryResource(did, resourceId)
	if err != nil {
		return types.DidDereferencing{}, err
	}
	if !isFound {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.NotFoundError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}
	jsonFragment, err := rs.didDocService.MarshallContentStream(&resource, dereferenceOptions.Accept)
	if err != nil {
		return types.DidDereferencing{}, err
	}
	contentStream := json.RawMessage(jsonFragment)

	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	return types.DidDereferencing{ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rs RequestService) dereferenceSecondary(did string, fragmentId string, didUrl string, dereferenceOptions types.DereferencingOption) (types.DidDereferencing, error) {
	didResolution, err := rs.Resolve(did, types.ResolutionOption(dereferenceOptions))
	if err != nil {
		return types.DidDereferencing{}, err
	}
	metadata := didResolution.Metadata
	dereferencingMetadata := types.DereferencingMetadata(didResolution.ResolutionMetadata)
	if dereferencingMetadata.ResolutionError != "" {
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	var protoContent protoiface.MessageV1
	if fragmentId != "" {
		protoContent = rs.didDocService.GetDIDFragment(fragmentId, didResolution.Did)
		metadata = types.TransformToFragmentMetadata(metadata)
	} else {
		protoContent = &didResolution.Did
	}

	if protoContent == nil {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.NotFoundError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	jsonFragment, err := rs.didDocService.MarshallContentStream(protoContent, dereferenceOptions.Accept)
	if err != nil {
		return types.DidDereferencing{}, err
	}
	contentStream := json.RawMessage(jsonFragment)

	return types.DidDereferencing{ContentStream: contentStream, Metadata: metadata, DereferencingMetadata: dereferencingMetadata}, nil
}

func (rs RequestService) ResolveMetadata(did string, metadata cheqdTypes.Metadata) (types.ResolutionDidDocMetadata, error) {
	if metadata.Resources == nil {
		return types.NewResolutionDidDocMetadata(did, metadata, []*resourceTypes.ResourceHeader{}), nil
	}
	resources, err := rs.ledgerService.QueryCollectionResources(did)
	if err != nil {
		return types.ResolutionDidDocMetadata{}, err
	}
	return types.NewResolutionDidDocMetadata(did, metadata, resources), nil
}

func createJsonResolution(didDoc string, metadata string, resolutionMetadata string) (string, error) {
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
		return "", err
	}

	return string(respJson), nil
}

func createJsonDereferencing(contentStream json.RawMessage, metadata string, dereferencingMetadata string) (string, error) {
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
		return "", err
	}

	return string(respJson), nil
}
