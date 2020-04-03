package transactiondomain

import (
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/types"
)

type TransactionSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	PublicKey types.PublicKey `json:"public_key"`
	Hash      types.Hash      `json:"hash"`
	Nonce     types.Nonce     `json:"nonce"`
	Fee       int64           `json:"fee"`
	GasLimit  uint64          `json:"gas_limit"`
	GasPrice  int64           `json:"gas_price"`
	Method    string          `json:"method"`
}

func (ts *TransactionSeq) ValidOwn() bool {
	return ts.PublicKey.Valid() &&
		ts.Hash.Valid() &&
		ts.Nonce.Valid()
}

func (ts *TransactionSeq) EqualOwn(m TransactionSeq) bool {
	return ts.PublicKey.Equal(m.PublicKey) &&
		ts.Hash.Equal(m.Hash)
}

func (ts *TransactionSeq) Valid() bool {
	return ts.DomainEntity.Valid() &&
		ts.Sequence.Valid() &&
		ts.ValidOwn()
}

func (ts *TransactionSeq) Equal(m TransactionSeq) bool {
	return ts.DomainEntity.Equal(*m.DomainEntity) &&
		ts.Sequence.Equal(*m.Sequence) &&
		ts.EqualOwn(m)
}
