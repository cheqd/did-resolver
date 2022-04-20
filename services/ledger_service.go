package services

import (
	"context"
	"errors"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type LedgerServiceI interface {
	QueryDIDDoc(did string) (cheqd.Did, cheqd.Metadata, bool, error)
	GetNamespaces() []string
}
type LedgerService struct {
	ledgers map[string]string // namespace -> url
}

func NewLedgerService() LedgerService {
	ls := LedgerService{}
	ls.ledgers = make(map[string]string)
	return ls
}

func (ls LedgerService) QueryDIDDoc(did string) (cheqd.Did, cheqd.Metadata, bool, error) {
	isFound := true
	serverAddr := ls.ledgers[getNamespace(did)]
	println(serverAddr)
	conn, err := openGRPCConnection(serverAddr)

	if err != nil {
		println("QueryDIDDoc: failed connection")
		isFound = false
		return cheqd.Did{}, cheqd.Metadata{}, isFound, err
	}
	println("QueryDIDDoc: successful connection")

	qc := cheqd.NewQueryClient(conn)
	defer conn.Close()

	didDocResponse, err := qc.Did(context.Background(), &cheqd.QueryGetDidRequest{Id: did})
	if err != nil {
		isFound = false
		return cheqd.Did{}, cheqd.Metadata{}, isFound, nil
	}
	println("QueryDIDDoc: received response")
	println(didDocResponse)
	return *didDocResponse.Did, *didDocResponse.Metadata, isFound, err
}

func (ls *LedgerService) RegisterLedger(namespace string, url string) error {
	println("RegisterLedger")

	if namespace == "" {
		println("Namespace cannot be empty")
		return errors.New("Namespace cannot be empty")
	}
	if url == "" {
		println("Ledger node url cannot be empty")
		return errors.New("Ledger node url cannot be empty")
	}
	ls.ledgers[namespace] = url

	println("RegisterLedger end")

	return nil
}

func openGRPCConnection(addr string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("ledgerTimeout"))
	defer cancel()

	conn, err = grpc.DialContext(ctx, addr, opts...)

	if err != nil {
		println("openGRPCConnection: context failed")
		println(err.Error())
		return nil, err
	}
	println("openGRPCConnection: opened")
	return conn, nil
}

func getNamespace(did string) string {
	_, namespace, _, _ := cheqdUtils.TrySplitDID(did)
	return namespace
}

func (ls LedgerService) GetNamespaces() []string {
	keys := make([]string, 0, len(ls.ledgers))
	for k, _ := range ls.ledgers {
		keys = append(keys, k)
	}
	return keys
}
