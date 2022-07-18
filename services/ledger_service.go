package services

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/credentials"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LedgerServiceI interface {
	QueryDIDDoc(did string) (cheqd.Did, cheqd.Metadata, bool, error)
	QueryResource(collectionDid string, resourceId string) (resource.Resource, bool, error)
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

func (ls LedgerService) QueryDIDDoc(did string) (cheqd.Did, cheqd.Metadata, bool, error) {
	_, namespace, _, _ := cheqdUtils.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[namespace]
	if !namespaceFound {
		return cheqd.Did{}, cheqd.Metadata{}, false, fmt.Errorf("namespace not supported: %s", namespace)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryDIDDoc: failed connection")
		return cheqd.Did{}, cheqd.Metadata{}, false, err
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying did doc: %s", did)
	client := cheqd.NewQueryClient(conn)
	didDocResponse, err := client.Did(context.Background(), &cheqd.QueryGetDidRequest{Id: did})
	if err != nil {
		return cheqd.Did{}, cheqd.Metadata{}, false, nil
	}

	return *didDocResponse.Did, *didDocResponse.Metadata, true, err
}

func (ls LedgerService) QueryResource(did string, resourceId string) (resource.Resource, bool, error) {
	_, namespace, collectionId, _ := cheqdUtils.TrySplitDID(did)
	serverAddr, namespaceFound := ls.ledgers[namespace]
	if !namespaceFound {
		return resource.Resource{}, false, fmt.Errorf("namespace not supported: %s", namespace)
	}

	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryResource: failed connection")
		return resource.Resource{}, false, err
	}

	defer mustCloseGRPCConnection(conn)

	log.Info().Msgf("Querying did resource: %s, %s", collectionId, resourceId)

	client := resource.NewQueryClient(conn)
	resourceResponse, err := client.Resource(context.Background(), &resource.QueryGetResourceRequest{CollectionId: collectionId, Id: resourceId})
	if err != nil {
		log.Info().Msgf("Resource not found %r", err.Error())

		return resource.Resource{}, false, nil
	}

	return *resourceResponse.Resource, true, err
}

func (ls *LedgerService) RegisterLedger(namespace string, url string) error {
	if namespace == "" {
		err := errors.New("namespace cannot be empty")
		log.Error().Err(err).Msg("RegisterLedger: failed")
		return err
	}

	if url == "" {
		return errors.New("ledger node url cannot be empty")
	}

	ls.ledgers[namespace] = url
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
		keys = append(keys, k)
	}
	return keys
}
