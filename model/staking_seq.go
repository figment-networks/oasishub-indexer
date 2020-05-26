package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type StakingSeq struct {
	*Model
	*Sequence

	TotalSupply         types.Quantity `json:"total_supply"`
	CommonPool          types.Quantity `json:"common_pool"`
	DebondingInterval   uint64         `json:"debonding_interval"`
	MinDelegationAmount types.Quantity `json:"min_delegation_amount"`
}

// - Methods
func (StakingSeq) TableName() string {
	return "staking_sequences"
}

func (ss *StakingSeq) Valid() bool {
	return ss.Sequence.Valid() &&
		ss.TotalSupply.Valid() &&
		ss.CommonPool.Valid()
}

func (ss *StakingSeq) Equal(m StakingSeq) bool {
	return ss.Sequence.Equal(*m.Sequence) &&
		ss.CommonPool.Equals(m.CommonPool) &&
		ss.TotalSupply.Equals(m.TotalSupply)
}
