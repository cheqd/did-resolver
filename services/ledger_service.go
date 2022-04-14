package services

import (
	"context"
	"errors"
	"flag"
	"strings"
	"time"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"google.golang.org/grpc"
)

type LedgerService struct {
	ledgers map[string]string // namespace -> url
}

func NewLedgerService() LedgerService {
	ls := LedgerService{}
	ls.ledgers = make(map[string]string)
	return ls
}

func (ls LedgerService) QueryDIDDoc(did string) (cheqd.Did, cheqd.Metadata, error) {
	serverAddr := ls.ledgers[getNamespace(did)]
	println(serverAddr)
	conn, err := openGRPCConnection(serverAddr)

	if err != nil {
		return cheqd.Did{}, cheqd.Metadata{}, err
	}

	qc := cheqd.NewQueryClient(conn)
	defer conn.Close()

	didDocResponse, err := qc.Did(context.Background(), &cheqd.QueryGetDidRequest{Id: did})

	println(didDocResponse)
	return *didDocResponse.Did, *didDocResponse.Metadata, err
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
	ls.ledgers[namespace] = *flag.String("grpc-server-address-"+namespace, url,
		"The target grpc server address in the format of host:port")

	println("RegisterLedger end")

	return nil
}

func openGRPCConnection(addr string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}
	// TODO: move to application setup
	// TODO: move timeouts to a config
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	return strings.SplitN(did, ":", 4)[2]
}
