package services

import (
	"errors"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
)

type LedgerService struct {
	ledgers map[string]string // namespace -> url
}

func (LedgerService) QueryDIDDoc(did string) (cheqd.Did, cheqd.Metadata, error) {
	return cheqd.Did{}, cheqd.Metadata{}, nil
}

func (ls *LedgerService) RegisterLedger(namespace string, url string) error {
	if !(namespace == "") {
		return errors.New("Namespace cannot be empty")
	}
	if !(url == "") {
		return errors.New("Ledger node url cannot be empty")
	}
	ls.ledgers[namespace] = url
	return nil
}
