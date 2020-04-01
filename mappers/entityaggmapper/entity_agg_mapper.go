package entityaggmapper

import (
	"github.com/figment-networks/oasishub/db/timescale/orm"
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/domain/entitydomain"
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/mappers/syncablemapper"
	"github.com/figment-networks/oasishub/types"
	"github.com/figment-networks/oasishub/utils/errors"
)

func FromPersistence(o orm.EntityAggModel) (*entitydomain.EntityAgg, errors.ApplicationError) {
	e := &entitydomain.EntityAgg{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: o.ID,
		}),
		Aggregate: commons.NewAggregate(commons.AggregateProps{
			StartedAtHeight: o.StartedAtHeight,
			StartedAt:       o.StartedAt,
		}),
		EntityUID: o.EntityUID,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("entity aggregate not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(r *entitydomain.EntityAgg) (*orm.EntityAggModel, errors.ApplicationError) {
	if !r.Valid() {
		return nil, errors.NewErrorFromMessage("entity aggregate not valid", errors.NotValid)
	}

	return &orm.EntityAggModel{
		EntityModel: orm.EntityModel{
			ID: r.ID,
		},
		AggregateModel: orm.AggregateModel{
			StartedAtHeight: r.GetStartedAtHeight(),
			StartedAt:       r.GetStartedAt(),
		},
		EntityUID: r.EntityUID,
	}, nil
}

func FromData(stateSyncable *syncabledomain.Syncable) ([]*entitydomain.EntityAgg, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var accounts []*entitydomain.EntityAgg
	for _, entity := range stateData.Data.Registry.Entities {
		acc := &entitydomain.EntityAgg{
			DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
			Aggregate: commons.NewAggregate(commons.AggregateProps{
				StartedAtHeight: stateSyncable.Height,
				StartedAt:       stateSyncable.Time,
			}),

			EntityUID: types.PublicKey(entity.Signature.PublicKey.String()),
		}

		if !acc.Valid() {
			return nil, errors.NewErrorFromMessage("entity aggregate not valid", errors.NotValid)
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}
