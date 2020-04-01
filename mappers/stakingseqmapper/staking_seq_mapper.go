package stakingseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/domain/stakingdomain"
	"github.com/figment-networks/oasishub-indexer/domain/syncabledomain"
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
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

func ToView(s *stakingdomain.StakingSeq) map[string]interface{} {
	return map[string]interface{}{
		"id":                    s.ID,
		"height":                s.Height,
		"time":                  s.Time,
		"chain_id":              s.ChainId,
		"total_supply":          s.TotalSupply.String(),
		"common_pool":           s.CommonPool.String(),
		"debonding_interval":    s.DebondingInterval,
		"min_delegation_amount": s.MinDelegationAmount.String(),
	}
}
