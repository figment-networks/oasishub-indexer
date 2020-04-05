package debondingdelegationseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/debondingdelegationseq"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func ToSequence(stateSyncable *syncable.Model) ([]debondingdelegationseq.Model, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var delegations []debondingdelegationseq.Model
	for validatorUID, delegationsMap := range stateData.Data.Staking.DebondingDelegations {
		for delegatorUID, infoArray := range delegationsMap {
			//TODO: Why is it array?
			info := infoArray[0]
			acc := debondingdelegationseq.Model{
				Sequence: &shared.Sequence{
					ChainId: stateSyncable.ChainId,
					Height:  stateSyncable.Height,
					Time:    stateSyncable.Time,
				},

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

func ToView(ts []*debondingdelegationseq.Model) []map[string]interface{} {
	var items []map[string]interface{}
	for _, t := range ts {
		i := map[string]interface{}{
			"id":       t.ID,
			"height":   t.Height,
			"time":     t.Time,
			"chain_id": t.ChainId,

			"validator_uid": t.ValidatorUID,
			"delegator_uid": t.DelegatorUID,
			"shares":        t.Shares.String(),
			"debond_end":    t.DebondEnd,
		}
		items = append(items, i)
	}
	return items
}
