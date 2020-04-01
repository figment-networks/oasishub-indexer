package stakingseqmapper

import (
	"github.com/figment-networks/oasishub/db/timescale/orm"
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/domain/stakingdomain"
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/mappers/syncablemapper"
	"github.com/figment-networks/oasishub/types"
	"github.com/figment-networks/oasishub/utils/errors"
)

func FromPersistence(o orm.StakingSeqModel) (*stakingdomain.StakingSeq, errors.ApplicationError) {
	e := &stakingdomain.StakingSeq{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: o.ID,
		}),
		Sequence: commons.NewSequence(commons.SequenceProps{
			ChainId: o.ChainId,
			Height:  o.Height,
			Time:    o.Time,
		}),

		TotalSupply:         o.TotalSupply,
		CommonPool:          o.CommonPool,
		DebondingInterval:   o.DebondingInterval,
		MinDelegationAmount: o.MinDelegationAmount,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("staking sequence not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(b *stakingdomain.StakingSeq) (*orm.StakingSeqModel, errors.ApplicationError) {
	if !b.Valid() {
		return nil, errors.NewErrorFromMessage("staking sequence not valid", errors.NotValid)
	}

	bs := &orm.StakingSeqModel{
		EntityModel: orm.EntityModel{ID: b.ID},
		SequenceModel: orm.SequenceModel{
			ChainId: b.ChainId,
			Height:  b.Height,
			Time:    b.Time,
		},

		TotalSupply:         b.TotalSupply,
		CommonPool:          b.CommonPool,
		DebondingInterval:   b.DebondingInterval,
		MinDelegationAmount: b.MinDelegationAmount,
	}
	return bs, nil
}

func FromData(stateSyncable syncabledomain.Syncable) (*stakingdomain.StakingSeq, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	e := &stakingdomain.StakingSeq{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
		Sequence: commons.NewSequence(commons.SequenceProps{
			ChainId: stateSyncable.ChainId,
			Height:  stateSyncable.Height,
			Time:    stateSyncable.Time,
		}),

		TotalSupply:         types.NewQuantity(stateData.Data.Staking.TotalSupply.ToBigInt()),
		CommonPool:          types.NewQuantity(stateData.Data.Staking.CommonPool.ToBigInt()),
		DebondingInterval:   uint64(stateData.Data.Staking.Parameters.DebondingInterval),
		MinDelegationAmount: types.NewQuantity(stateData.Data.Staking.Parameters.MinDelegationAmount.ToBigInt()),
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("staking sequence not valid", errors.NotValid)
	}

	return e, nil
}
