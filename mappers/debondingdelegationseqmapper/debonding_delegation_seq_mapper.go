package debondingdelegationseqmapper

import (
	"github.com/figment-networks/oasishub/db/timescale/orm"
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/domain/delegationdomain"
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/mappers/syncablemapper"
	"github.com/figment-networks/oasishub/types"
	"github.com/figment-networks/oasishub/utils/errors"
)

func FromPersistence(o orm.DebondingDelegationSeqModel) (*delegationdomain.DebondingDelegationSeq, errors.ApplicationError) {
	e := &delegationdomain.DebondingDelegationSeq{
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
		Shares:       o.Shares,
		DebondEnd:    o.DebondEnd,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("debonding delegation sequence not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(r *delegationdomain.DebondingDelegationSeq) (*orm.DebondingDelegationSeqModel, errors.ApplicationError) {
	if !r.Valid() {
		return nil, errors.NewErrorFromMessage("debonding delegation sequence not valid", errors.NotValid)
	}

	return &orm.DebondingDelegationSeqModel{
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
		Shares:       r.Shares,
		DebondEnd:    r.DebondEnd,
	}, nil
}

func FromData(stateSyncable *syncabledomain.Syncable) ([]*delegationdomain.DebondingDelegationSeq, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var delegations []*delegationdomain.DebondingDelegationSeq
	for validatorUID, delegationsMap := range stateData.Data.Staking.DebondingDelegations {
		for delegatorUID, infoArray := range delegationsMap {
			//TODO: Why is it array?
			info := infoArray[0]
			acc := &delegationdomain.DebondingDelegationSeq{
				DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
				Sequence: commons.NewSequence(commons.SequenceProps{
					ChainId: stateSyncable.ChainId,
					Height:  stateSyncable.Height,
					Time:    stateSyncable.Time,
				}),

				ValidatorUID: types.PublicKey(validatorUID.String()),
				DelegatorUID: types.PublicKey(delegatorUID.String()),
				Shares:       types.NewQuantity(info.Shares.ToBigInt()),
				DebondEnd:    int64(info.DebondEndTime),
			}

			if !acc.Valid() {
				return nil, errors.NewErrorFromMessage("debonding delegation sequence not valid", errors.NotValid)
			}

			delegations = append(delegations, acc)
		}
	}
	return delegations, nil
}
