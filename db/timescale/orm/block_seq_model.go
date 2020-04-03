package orm

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type BlockSeqModel struct {
	EntityModel
	SequenceModel

	Hash              types.Hash
	ProposerEntityUID types.PublicKey
	AppVersion        int64
	BlockVersion      int64
	TransactionsCount types.Count
}

func (BlockSeqModel) TableName() string {
	return "block_sequences"
}
