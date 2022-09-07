package services

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"time"

	"google.golang.org/grpc/credentials"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DELIMITER = ":"
)

type LedgerServiceI interface {
	QueryDIDDoc(did string) (*cheqd.Did, *cheqd.Metadata, *types.IdentityError)
	QueryResource(collectionDid string, resourceId string) (*resource.Resource, *types.IdentityError)
	QueryCollectionResources(did string) ([]*resource.ResourceHeader, *types.IdentityError)
	GetNamespaces() []string
}

type LedgerService struct {
	ledgers           map[string]string // namespace -> url
	connectionTimeout time.Duration
	useTls            bool
}

func NewLedgerService(connectionTimeout time.Duration, useTls bool) LedgerService {
	ls := LedgerService{
		connectionTimeout: connectionTimeout,
		useTls:            useTls,
	}
	ls.ledgers = make(map[string]string)
	return ls
}

func (ls LedgerService) QueryDIDDoc(did string) (*cheqd.Did, *cheqd.Metadata, *types.IdentityError) {
	method, namespace, _, _ := cheqdUtils.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, nil, types.NewInvalidDIDError(did, types.JSON, nil, false)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryDIDDoc: failed connection")
		return nil, nil, types.NewInternalError(did, types.JSON, err, false)
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying did doc: %s", did)
	client := cheqd.NewQueryClient(conn)
	didDocResponse, err := client.Did(context.Background(), &cheqd.QueryGetDidRequest{Id: did})
	if err != nil {
		return nil, nil, types.NewNotFoundError(did, types.JSON, err, false)
	}

	return didDocResponse.Did, didDocResponse.Metadata, nil
}

func (ls LedgerService) QueryResource(did string, resourceId string) (*resource.Resource, *types.IdentityError) {
	method, namespace, collectionId, _ := cheqdUtils.TrySplitDID(did)
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

	log.Info().Msgf("Querying did resource: %s, %s", collectionId, resourceId)

	client := resource.NewQueryClient(conn)
	resourceResponse, err := client.Resource(context.Background(), &resource.QueryGetResourceRequest{CollectionId: collectionId, Id: resourceId})
	if err != nil {
		log.Info().Msgf("Resource not found %s", err.Error())
		return nil, types.NewNotFoundError(did, types.JSON, err, true)
	}

	return resourceResponse.Resource, nil
}

func (ls LedgerService) QueryCollectionResources(did string) ([]*resource.ResourceHeader, *types.IdentityError) {
	method, namespace, collectionId, _ := cheqdUtils.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[method+DELIMITER+namespace]
	if !namespaceFound {
		return nil, types.NewInvalidDIDError(did, types.JSON, nil, false)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryResource: failed connection")
		return nil, types.NewInternalError(did, types.JSON, err, false)
	}

	log.Info().Msgf("Querying did resources: %s", did)

	client := resource.NewQueryClient(conn)
	resourceResponse, err := client.CollectionResources(context.Background(), &resource.QueryGetCollectionResourcesRequest{CollectionId: collectionId})
	if err != nil {
		return nil, types.NewNotFoundError(did, types.JSON, err, false)
	}

	return resourceResponse.Resources, nil
}

func (ls *LedgerService) RegisterLedger(method string, namespace string, url string) error {
	if namespace == "" || method == "" {
		err := errors.New("namespace and method cannot be empty")
		log.Error().Err(err).Msg("RegisterLedger: failed")
		return err
	}

	if url == "" {
		return errors.New("ledger node url cannot be empty")
	}

	ls.ledgers[method+DELIMITER+namespace] = url
	return nil
}

func (ls LedgerService) openGRPCConnection(addr string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	if ls.useTls {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), ls.connectionTimeout)
	defer cancel()

	conn, err = grpc.DialContext(ctx, addr, opts...)

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
