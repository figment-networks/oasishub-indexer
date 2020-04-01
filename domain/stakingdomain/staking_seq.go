package stakingdomain

import (
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/types"
)

type StakingSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	TotalSupply         types.Quantity
	CommonPool          types.Quantity
	DebondingInterval   uint64
	MinDelegationAmount types.Quantity
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
