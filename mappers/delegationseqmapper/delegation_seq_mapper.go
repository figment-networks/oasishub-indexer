package delegationseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/domain/delegationdomain"
	"github.com/figment-networks/oasishub-indexer/domain/syncabledomain"
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func FromPersistence(o orm.DelegationSeqModel) (*delegationdomain.DelegationSeq, errors.ApplicationError) {
	e := &delegationdomain.DelegationSeq{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: o.ID,
		}),
		Sequence: commons.NewSequence(commons.SequenceProps{
			ChainId: o.ChainId,
			Height:  o.Height,
			Time:    o.Time,
		}),
		ValidatorUID: o.ValidatorUID,
		DelegatorUID: o.DelegatorUID,
		Shares: o.Shares,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("delegation sequence not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(r *delegationdomain.DelegationSeq) (*orm.DelegationSeqModel, errors.ApplicationError) {
	if !r.Valid() {
		return nil, errors.NewErrorFromMessage("delegation sequence not valid", errors.NotValid)
	}

	return &orm.DelegationSeqModel{
		EntityModel: orm.EntityModel{
			ID: r.ID,
		},
		SequenceModel: orm.SequenceModel{
			ChainId: r.ChainId,
			Height:  r.Height,
			Time:    r.Time,
		},
		ValidatorUID: r.ValidatorUID,
		DelegatorUID: r.DelegatorUID,
		Shares: r.Shares,
	}, nil
}

func FromData(stateSyncable *syncabledomain.Syncable) ([]*delegationdomain.DelegationSeq, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var delegations []*delegationdomain.DelegationSeq
	for validatorUID, delegationsMap := range stateData.Data.Staking.Delegations {
		for delegatorUID, info := range delegationsMap {
			acc := &delegationdomain.DelegationSeq{
				DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
				Sequence: commons.NewSequence(commons.SequenceProps{
					ChainId: stateSyncable.ChainId,
					Height:  stateSyncable.Height,
					Time:    stateSyncable.Time,
				}),

				ValidatorUID: types.PublicKey(validatorUID.String()),
				DelegatorUID: types.PublicKey(delegatorUID.String()),
				Shares:       types.NewQuantity(info.Shares.ToBigInt()),
			}

			if !acc.Valid() {
				return nil, errors.NewErrorFromMessage("delegation sequence not valid", errors.NotValid)
			}

			delegations = append(delegations, acc)
		}
	}
	return delegations, nil
}