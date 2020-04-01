package orm

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type TransactionSeqModel struct {
	EntityModel
	SequenceModel

	PublicKey types.PublicKey
	Hash      types.Hash
	Nonce     types.Nonce
	Fee       int64
	GasLimit  uint64
	GasPrice  int64
	Method    string
}

func (TransactionSeqModel) TableName() string {
	return "transaction_sequences"
}
