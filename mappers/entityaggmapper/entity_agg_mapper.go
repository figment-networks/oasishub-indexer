package entityaggmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/entityagg"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func ToAggregate(stateSyncable *syncable.Model) ([]entityagg.Model, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var accounts []entityagg.Model
	for _, entity := range stateData.Data.Registry.Entities {
		acc := entityagg.Model{
			Aggregate: &shared.Aggregate{
				StartedAtHeight: stateSyncable.Height,
				StartedAt:       stateSyncable.Time,
			},

			EntityUID: types.PublicKey(entity.Signature.PublicKey.String()),
		}

		if !acc.Valid() {
			return nil, errors.NewErrorFromMessage("entity aggregate not valid", errors.NotValid)
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}
