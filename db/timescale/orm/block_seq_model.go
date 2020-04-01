package orm

import (
	"github.com/figment-networks/oasishub/types"
)

type BlockSeqModel struct {
	EntityModel
	SequenceModel

	AppVersion        int64
	BlockVersion      int64
	TransactionsCount types.Count
}

func (BlockSeqModel) TableName() string {
	return "block_sequences"
}
