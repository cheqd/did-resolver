package services

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"

	"google.golang.org/grpc/credentials"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/types"
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
	ledgers map[string]types.Network // namespace -> endpoint with configs
}

func NewLedgerService() LedgerService {
	ls := LedgerService{}
	ls.ledgers = make(map[string]types.Network)

	return ls
}

func (ls LedgerService) QueryDIDDoc(did string, version string) (*didTypes.DidDocWithMetadata, *types.IdentityError) {
	method, namespace, _, _ := types.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, false)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryDIDDoc: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying DIDDoc: %s", did)
	client := didTypes.NewQueryClient(conn)

	if version == "" {
		didDocResponse, err := client.DidDoc(context.Background(), &didTypes.QueryDidDocRequest{Id: did})
		if err != nil {
			return nil, types.NewNotFoundError(did, types.JSON, err, false)
		}

		return didDocResponse.Value, nil
	} else {
		didDocResponse, err := client.DidDocVersion(context.Background(), &didTypes.QueryDidDocVersionRequest{Id: did, Version: version})
		if err != nil {
			return nil, types.NewNotFoundError(did, types.JSON, err, false)
		}

		return didDocResponse.Value, nil
	}
}

func (ls LedgerService) QueryAllDidDocVersionsMetadata(did string) ([]*didTypes.Metadata, *types.IdentityError) {
	method, namespace, _, _ := types.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, false)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryAllDidDocVersionsMetadata: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}
	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying all DIDDoc versions metadata: %s", did)
	client := didTypes.NewQueryClient(conn)

	response, err := client.AllDidDocVersionsMetadata(context.Background(), &didTypes.QueryAllDidDocVersionsMetadataRequest{Id: did})
	if err != nil {
		return nil, types.NewNotFoundError(did, types.JSON, err, false)
	}

	return response.Versions, nil
}

func (ls LedgerService) QueryResource(did string, resourceId string) (*resourceTypes.ResourceWithMetadata, *types.IdentityError) {
	method, namespace, collectionId, _ := types.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, true)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryResource: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, true)
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying DID resource: %s, %s", collectionId, resourceId)

	client := resourceTypes.NewQueryClient(conn)
	resourceResponse, err := client.Resource(context.Background(), &resourceTypes.QueryResourceRequest{CollectionId: collectionId, Id: resourceId})
	if err != nil {
		log.Info().Msgf("Resource not found %s", err.Error())
		return nil, types.NewNotFoundError(did, types.JSON, err, true)
	}

	return resourceResponse.Resource, nil
}

func (ls LedgerService) QueryCollectionResources(did string) ([]*resourceTypes.Metadata, *types.IdentityError) {
	method, namespace, collectionId, _ := types.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, false)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryResource: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}

	log.Info().Msgf("Querying DID resources: %s", did)

	client := resourceTypes.NewQueryClient(conn)
	resourceResponse, err := client.CollectionResources(context.Background(), &resourceTypes.QueryCollectionResourcesRequest{CollectionId: collectionId})
	if err != nil {
		return nil, types.NewNotFoundError(did, types.JSON, err, false)
	}

	return resourceResponse.Resources, nil
}

func (ls *LedgerService) RegisterLedger(method string, endpoint types.Network) error {
	if endpoint.Namespace == "" || method == "" {
		err := errors.New("namespace and method cannot be empty")
		log.Error().Err(err).Msg("RegisterLedger: failed")
		return err
	}

	if endpoint.Endpoint == "" {
		return errors.New("ledger node URL cannot be empty")
	}

	ls.ledgers[method+DELIMITER+endpoint.Namespace] = endpoint

	return nil
}

func (ls LedgerService) openGRPCConnection(endpoint types.Network) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	if endpoint.UseTls {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), endpoint.Timeout)
	defer cancel()

	conn, err = grpc.DialContext(ctx, endpoint.Endpoint, opts...)

	if err != nil {
		log.Error().Err(err).Msgf("openGRPCConnection: context failed")
		return nil, err
	}

	log.Info().Msg("openGRPCConnection: opened")

	return conn, nil
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
