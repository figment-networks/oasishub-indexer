package stakingdomain

import (
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/types"
)

type StakingSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	TotalSupply         types.Quantity `json:"total_supply"`
	CommonPool          types.Quantity `json:"common_pool"`
	DebondingInterval   uint64         `json:"debonding_interval"`
	MinDelegationAmount types.Quantity `json:"min_delegation_amount"`
}

func (ss *StakingSeq) ValidOwn() bool {
	return ss.TotalSupply.Valid() &&
		ss.CommonPool.Valid()
}

func (ss *StakingSeq) EqualOwn(m StakingSeq) bool {
	return true
}

func (ss *StakingSeq) Valid() bool {
	return ss.DomainEntity.Valid() &&
		ss.Sequence.Valid() &&
		ss.ValidOwn()
}

func (ss *StakingSeq) Equal(m StakingSeq) bool {
	return ss.DomainEntity.Equal(*m.DomainEntity) &&
		ss.Sequence.Equal(*m.Sequence) &&
		ss.EqualOwn(m)
}
