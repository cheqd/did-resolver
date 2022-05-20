package services

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"google.golang.org/grpc/credentials"
	"time"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LedgerServiceI interface {
	QueryDIDDoc(did string) (cheqd.Did, cheqd.Metadata, bool, error)
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

	log.Info().Msgf("Connecting to the ledger: %s", serverAddr)
	conn, err := ls.openGRPCConnection(serverAddr)
	if err != nil {
		log.Error().Err(err).Msg("QueryDIDDoc: failed connection")
		return cheqd.Did{}, cheqd.Metadata{}, false, err
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Panic().Err(err).Msg("QueryDIDDoc: failed to close connection")
			panic(err)
		}
	}(conn)

	log.Info().Msgf("Querying did doc: %s", did)
	client := cheqd.NewQueryClient(conn)
	didDocResponse, err := client.Did(context.Background(), &cheqd.QueryGetDidRequest{Id: did})
	if err != nil {
		return cheqd.Did{}, cheqd.Metadata{}, false, nil
	}

	return *didDocResponse.Did, *didDocResponse.Metadata, true, err
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

func (ls LedgerService) GetNamespaces() []string {
	keys := make([]string, 0, len(ls.ledgers))
	for k := range ls.ledgers {
		keys = append(keys, k)
	}
	return keys
}
