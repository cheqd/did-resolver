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

func (rs RequestService) ProcessDIDRequest(didUrl string, resolutionOptions types.ResolutionOption) (string, int) {
	var result string
	var statusCode int
	if utils.IsDidUrl(didUrl) {
		log.Trace().Msgf("Dereferencing %s", didUrl)
		result, statusCode = rs.prepareDereferencingResult(didUrl, types.DereferencingOption(resolutionOptions))
	} else {
		log.Trace().Msgf("Resolving %s", didUrl)
		result, statusCode = rs.prepareResolutionResult(didUrl, resolutionOptions)
	}
	return result, statusCode
}

func (rs RequestService) prepareResolutionResult(did string, resolutionOptions types.ResolutionOption) (string, int) {
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

func (rs RequestService) prepareDereferencingResult(did string, dereferencingOptions types.DereferencingOption) (string, int) {
	log.Info().Msgf("Dereferencing %s", did)

	didDereferencing := rs.Dereference(did, dereferencingOptions)

	dereferencingMetadata, mErr1 := json.Marshal(didDereferencing.DereferencingMetadata)
	metadata, mErr2 := json.Marshal(didDereferencing.Metadata)
	if mErr1 != nil || mErr2 != nil {
		log.Error().Errs("errors", []error{mErr1, mErr2}).Msg("Marshalling error")
		return createJsonDereferencingInternalError([]byte{})
	}

	if didDereferencing.DereferencingMetadata.ResolutionError != "" {
		didDereferencing.ContentStream = nil
		metadata = []byte{}
	}

	result, err := createJsonDereferencing(didDereferencing.ContentStream, string(metadata), string(dereferencingMetadata))
	if err != nil {
		log.Error().Err(err).Msg("Marshalling error")
		return createJsonDereferencingInternalError(dereferencingMetadata)
	}

	return result, didDereferencing.DereferencingMetadata.ResolutionError.GetStatusCode()
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

	resolvedMetadata, err := rs.ResolveMetadata(did, metadata)
	if err != nil {
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
func (rs RequestService) Dereference(didUrl string, dereferenceOptions types.DereferencingOption) types.DidDereferencing {
	did, path, query, fragmentId, err := cheqdUtils.TrySplitDIDUrl(didUrl)
	log.Info().Msgf("did: %s, path: %s, query: %s, fragmentId: %s", did, path, query, fragmentId)

	if !dereferenceOptions.Accept.IsSupported() {
		return types.DidDereferencing{DereferencingMetadata: types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)}
	}

	if err != nil || !cheqdUtils.IsValidDIDUrl(didUrl, "", []string{}) {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.InvalidDIDUrlError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	// TODO: implement
	var didDereferencing types.DidDereferencing
	if query != "" {
		didDereferencing, err = rs.dereferenceService(did, query, didUrl, dereferenceOptions)
	} else if path != "" {
		didDereferencing, err = rs.dereferencePrimary(path, did, didUrl, dereferenceOptions)
	} else {
		didDereferencing, err = rs.dereferenceSecondary(did, fragmentId, didUrl, dereferenceOptions)
	}

	if err != nil {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.InternalError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	return didDereferencing
}

func (rs RequestService) dereferenceService(did string, query string, didUrl string, dereferenceOptions types.DereferencingOption) (types.DidDereferencing, error) {
	didResolution := rs.Resolve(did, types.ResolutionOption(dereferenceOptions))

	dereferencingMetadata := types.DereferencingMetadata(didResolution.ResolutionMetadata)
	if dereferencingMetadata.ResolutionError != "" {
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	service := rs.didDocService.GetDIDQuery(query, didResolution.Did)

	if service == nil {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.NotFoundError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	service.ServiceEndpoint = CreatServiceEndpoint(didUrl, service.ServiceEndpoint)
	metadata := types.TransformToFragmentMetadata(didResolution.Metadata)

	jsonFragment, err := rs.didDocService.MarshallContentStream(service, dereferenceOptions.Accept)
	if err != nil {
		return types.DidDereferencing{}, err
	}
	contentStream := json.RawMessage(jsonFragment)

	dereferencingMetadata = types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	return types.DidDereferencing{ContentStream: contentStream, Metadata: metadata, DereferencingMetadata: dereferencingMetadata}, nil
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
	didResolution := rs.Resolve(did, types.ResolutionOption(dereferenceOptions))

	dereferencingMetadata := types.DereferencingMetadata(didResolution.ResolutionMetadata)
	if dereferencingMetadata.ResolutionError != "" {
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
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

func createJsonDereferencingInternalError(dereferencingMetadata []byte) (string, int) {
	result, mErr := createJsonDereferencing(nil, "", string(dereferencingMetadata))
	if mErr != nil {
		return "", types.InternalError.GetStatusCode()
	}
	return result, types.InternalError.GetStatusCode()
}

func createJsonResolutionInternalError(resolutionMetadata []byte) (string, int) {
	result, mErr := createJsonResolution("", "", string(resolutionMetadata))
	if mErr != nil {
		return "", types.InternalError.GetStatusCode()
	}
	return result, types.InternalError.GetStatusCode()
}

func CreatServiceEndpoint(didUrl string, inputServiceEndpoint string) (outputServiceEndpoint string) {
	_, path, query, fragment, _ := cheqdUtils.TrySplitDIDUrl(didUrl)
	outputServiceEndpoint = inputServiceEndpoint
	if path != "" {
		outputServiceEndpoint += "/" + path
		if query != "" {
			outputServiceEndpoint += "?" + strings.Split(query, "=")[1]
		}
	}
	if fragment != "" {
		outputServiceEndpoint += "#" + fragment
	}
	return outputServiceEndpoint
}
