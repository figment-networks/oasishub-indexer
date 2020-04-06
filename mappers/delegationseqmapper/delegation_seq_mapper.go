package delegationseqmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/delegationseq"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func ToSequence(stateSyncable *syncable.Model) ([]delegationseq.Model, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var delegations []delegationseq.Model
	for validatorUID, delegationsMap := range stateData.Data.Staking.Delegations {
		for delegatorUID, info := range delegationsMap {
			acc := delegationseq.Model{
				Sequence: &shared.Sequence{
					ChainId: stateSyncable.ChainId,
					Height:  stateSyncable.Height,
					Time:    stateSyncable.Time,
				},

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

type ListView struct {
	Items []delegationseq.Model `json:"items"`
}

func ToListView(ms []delegationseq.Model) *ListView {
	return &ListView{
		Items: ms,
	}
}
