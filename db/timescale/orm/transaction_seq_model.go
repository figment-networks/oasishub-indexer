package orm

import (
	"github.com/figment-networks/oasishub/types"
)

type TransactionSeqModel struct {
	EntityModel
	SequenceModel

	Hash     types.Hash
	Fee      int64
	GasLimit uint64
	GasPrice int64
	Method   string
}

func (TransactionSeqModel) TableName() string {
	return "transaction_sequences"
}