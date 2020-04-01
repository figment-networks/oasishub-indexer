package orm

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type StakingSeqModel struct {
	EntityModel
	SequenceModel

	// Associations

	// Indexes
	TotalSupply         types.Quantity
	CommonPool          types.Quantity
	DebondingInterval   uint64
	MinDelegationAmount types.Quantity
}

func (StakingSeqModel) TableName() string {
	return "staking_sequences"
}
