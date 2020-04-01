package transactiondomain

import (
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/types"
)

type TransactionSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	Hash     types.Hash
	Fee      int64
	GasLimit uint64
	GasPrice int64
	Method   string
}

func (ts *TransactionSeq) ValidOwn() bool {
	return ts.Hash.Valid()
}

func (ts *TransactionSeq) EqualOwn(m TransactionSeq) bool {
	return ts.Hash.Equal(m.Hash)
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
