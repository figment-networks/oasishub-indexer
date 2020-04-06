package stakingseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/stakingseq"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func ToSequence(stateSyncable syncable.Model) (*stakingseq.Model, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	e := &stakingseq.Model{
		Sequence: &shared.Sequence{
			ChainId: stateSyncable.ChainId,
			Height:  stateSyncable.Height,
			Time:    stateSyncable.Time,
		},

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

type DetailsView struct {
	*stakingseq.Model
}

func ToDetailsView(s *stakingseq.Model) *DetailsView {
	return &DetailsView{
		Model: s,
	}
}
