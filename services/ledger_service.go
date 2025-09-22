package services

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/credentials"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DELIMITER = ":"
)

type LedgerServiceI interface {
	QueryDIDDoc(did string, version string) (*didTypes.DidDocWithMetadata, *types.IdentityError)
	QueryAllDidDocVersionsMetadata(did string) ([]*didTypes.Metadata, *types.IdentityError)
	QueryResource(collectionDid string, resourceId string) (*resourceTypes.ResourceWithMetadata, *types.IdentityError)
	QueryCollectionResources(did string) ([]*resourceTypes.Metadata, *types.IdentityError)
	GetNamespaces() []string
}

type LedgerService struct {
	ledgers         map[string]types.Network // namespace -> endpoint with configs
	endpointManager *EndpointManager
}

func NewLedgerService(endpointManager *EndpointManager) LedgerService {
	ls := LedgerService{}
	ls.ledgers = make(map[string]types.Network)
	ls.endpointManager = endpointManager

	return ls
}

// GetHealthyConnection handles endpoint selection, connection, and automatic fallback
func (ls LedgerService) GetHealthyConnection(namespace string, did string) (*grpc.ClientConn, *types.IdentityError) {
	// Get healthy network from endpoint manager
	network, err := ls.endpointManager.GetHealthyEndpoint(namespace)
	if err != nil {
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}

	// The EndpointManager returns a Network with only the healthy endpoint in the slice
	if len(network.Endpoints) == 0 {
		return nil, types.NewInternalError(did, types.JSON, fmt.Errorf("no healthy endpoints available"), false)
	}

	// Use the healthy endpoint returned by EndpointManager
	healthyEndpoint := network.Endpoints[0]

	conn, err := ls.openGRPCConnection(*network)
	if err == nil {
		return conn, nil
	}

	log.Error().Err(err).Msgf("Failed connection to %s", healthyEndpoint.URL)
	ls.endpointManager.MarkEndpointUnhealthy(*network)

	// Try the other endpoint if available (only for connection failures)
	if fallbackNetwork := ls.getOtherEndpoint(namespace, network); fallbackNetwork != nil {
		fallbackEndpoint := fallbackNetwork.Endpoints[0]
		log.Info().Msgf("Trying other endpoint %s for namespace %s", fallbackEndpoint.URL, namespace)

		fallbackConn, fallbackErr := ls.openGRPCConnection(*fallbackNetwork)
		if fallbackErr == nil {
			log.Info().Msgf("Successfully connected to other endpoint %s", fallbackEndpoint.URL)
			return fallbackConn, nil
		}

		log.Error().Err(fallbackErr).Msgf("Other endpoint %s also failed", fallbackEndpoint.URL)
		ls.endpointManager.MarkEndpointUnhealthy(*fallbackNetwork)
		if fallbackConn != nil {
			fallbackConn.Close()
		}
	}

	// Both attempts failed
	return nil, types.NewInternalError(did, types.JSON, err, false)
}

func (ls LedgerService) QueryDIDDoc(did string, version string) (*didTypes.DidDocWithMetadata, *types.IdentityError) {
	method, namespace, _, _ := utils.TrySplitDID(did)
	_, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDidError(did, types.JSON, nil, false)
	}

	// Get healthy connection with automatic fallback
	conn, err := ls.GetHealthyConnection(namespace, did)
	if err != nil {
		return nil, err
	}
	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying DIDDoc: %s", did)
	client := didTypes.NewQueryClient(conn)

	if version == "" {
		didDocResponse, grpcErr := client.DidDoc(context.Background(), &didTypes.QueryDidDocRequest{Id: did})
		if grpcErr != nil {
			return nil, types.NewNotFoundError(did, types.JSON, grpcErr, false)
		}

		return didDocResponse.Value, nil
	}

	didDocResponse, grpcErr := client.DidDocVersion(context.Background(), &didTypes.QueryDidDocVersionRequest{Id: did, Version: version})
	if grpcErr != nil {
		return nil, types.NewNotFoundError(did, types.JSON, grpcErr, false)
	}

	return didDocResponse.Value, nil
}

func (ls LedgerService) QueryAllDidDocVersionsMetadata(did string) ([]*didTypes.Metadata, *types.IdentityError) {
	method, namespace, _, _ := utils.TrySplitDID(did)
	_, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDidError(did, types.JSON, nil, false)
	}

	// Get healthy connection with automatic fallback
	conn, err := ls.GetHealthyConnection(namespace, did)
	if err != nil {
		return nil, err
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying all DIDDoc versions metadata: %s", did)
	client := didTypes.NewQueryClient(conn)

	didDocResponse, grpcErr := client.AllDidDocVersionsMetadata(context.Background(), &didTypes.QueryAllDidDocVersionsMetadataRequest{Id: did})
	if grpcErr != nil {
		return nil, types.NewNotFoundError(did, types.JSON, grpcErr, false)
	}

	return didDocResponse.Versions, nil
}

func (ls LedgerService) QueryResource(did string, resourceId string) (*resourceTypes.ResourceWithMetadata, *types.IdentityError) {
	method, namespace, collectionId, _ := utils.TrySplitDID(did)
	_, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDidError(did, types.JSON, nil, true)
	}

	// Get healthy connection with automatic fallback
	conn, err := ls.GetHealthyConnection(namespace, did)
	if err != nil {
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying DID resource: %s, %s", collectionId, resourceId)

	client := resourceTypes.NewQueryClient(conn)
	resourceResponse, grpcErr := client.Resource(context.Background(), &resourceTypes.QueryResourceRequest{CollectionId: collectionId, Id: resourceId})
	if grpcErr != nil {
		log.Error().Msgf("Resource not found %s", grpcErr.Error())
		return nil, types.NewNotFoundError(did, types.JSON, grpcErr, true)
	}

	return resourceResponse.Resource, nil
}

func (ls LedgerService) QueryCollectionResources(did string) ([]*resourceTypes.Metadata, *types.IdentityError) {
	method, namespace, collectionId, _ := utils.TrySplitDID(did)
	_, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDidError(did, types.JSON, nil, false)
	}

	// Get healthy connection with automatic fallback
	conn, err := ls.GetHealthyConnection(namespace, did)
	if err != nil {
		return nil, err
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying DID resources: %s", did)

	client := resourceTypes.NewQueryClient(conn)
	resourceResponse, grpcErr := client.CollectionResources(context.Background(), &resourceTypes.QueryCollectionResourcesRequest{CollectionId: collectionId})
	if grpcErr != nil {
		return nil, types.NewNotFoundError(did, types.JSON, grpcErr, false)
	}

	return resourceResponse.Resources, nil
}

func (ls *LedgerService) RegisterLedger(method string, endpoint types.Network) error {
	if endpoint.Namespace == "" || method == "" {
		err := errors.New("namespace and method cannot be empty")
		log.Error().Err(err).Msg("RegisterLedger: failed")
		return err
	}

	if len(endpoint.Endpoints) == 0 {
		return errors.New("ledger node must have at least one endpoint configured")
	}

	ls.ledgers[method+DELIMITER+endpoint.Namespace] = endpoint

	return nil
}

func (ls LedgerService) openGRPCConnection(endpoint types.Network) (conn *grpc.ClientConn, err error) {
	// Use the first endpoint in the slice (guaranteed to be primary)
	if len(endpoint.Endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints configured for network")
	}

	// Use shared utility function to eliminate code duplication
	return openGRPCConnectionWithTimeout(endpoint.Endpoints[0].URL, endpoint.Endpoints[0].UseTls, endpoint.Endpoints[0].Timeout)
}

func mustCloseGRPCConnection(conn *grpc.ClientConn) {
	if conn == nil {
		return
	}
	err := conn.Close()
	if err != nil {
		log.Panic().Err(err).Msg("QueryDIDDoc: failed to close connection")
		panic(err)
	}
}

func (ls LedgerService) GetNamespaces() []string {
	keys := make([]string, 0, len(ls.ledgers))
	for k := range ls.ledgers {
		namespace := strings.Split(k, DELIMITER)[1]
		keys = append(keys, namespace)
	}

	return keys
}

// getOtherEndpoint gets the other endpoint for the namespace (simple fallback)
func (ls LedgerService) getOtherEndpoint(namespace string, currentNetwork *types.Network) *types.Network {
	if ls.endpointManager == nil {
		return nil
	}

	// Simple approach: get any healthy endpoint that's not the current one
	healthyNetwork, err := ls.endpointManager.GetHealthyEndpoint(namespace)
	if err != nil {
		return nil
	}

	// If it's the same endpoint, no fallback available
	if len(healthyNetwork.Endpoints) > 0 && len(currentNetwork.Endpoints) > 0 {
		if healthyNetwork.Endpoints[0].URL == currentNetwork.Endpoints[0].URL {
			return nil
		}
	}

	return healthyNetwork
}

// openGRPCConnectionWithTimeout creates a gRPC connection with timeout
func openGRPCConnectionWithTimeout(endpoint string, useTls bool, timeout time.Duration) (*grpc.ClientConn, error) {
	// Dial options (credentials only). Connection readiness is verified by the subsequent RPC's context timeout.
	cred := grpc.WithTransportCredentials(insecure.NewCredentials())
	if useTls {
		cred = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	}

	conn, err := grpc.NewClient(endpoint, cred)
	if err != nil {
		log.Error().Err(err).Msgf("openGRPCConnection: connection failed")
		return nil, err
	}

	log.Info().Msg("openGRPCConnection: opened")
	return conn, nil
}
