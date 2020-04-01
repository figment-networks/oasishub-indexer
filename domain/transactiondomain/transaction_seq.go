package transactiondomain

import (
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/types"
)

type TransactionSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	PublicKey types.PublicKey
	Hash      types.Hash
	Nonce     types.Nonce
	Fee       int64
	GasLimit  uint64
	GasPrice  int64
	Method    string
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
