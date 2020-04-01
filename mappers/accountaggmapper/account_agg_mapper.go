package accountaggmapper

import (
	"github.com/figment-networks/oasishub-indexer/db/timescale/orm"
	"github.com/figment-networks/oasishub-indexer/domain/accountdomain"
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/domain/syncabledomain"
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func FromPersistence(o orm.AccountAggModel) (*accountdomain.AccountAgg, errors.ApplicationError) {
	e := &accountdomain.AccountAgg{
		DomainEntity: commons.NewDomainEntity(commons.EntityProps{
			ID: o.ID,
		}),
		Aggregate: commons.NewAggregate(commons.AggregateProps{
			StartedAtHeight: o.StartedAtHeight,
			StartedAt:       o.StartedAt,
		}),
		PublicKey:                      o.PublicKey,
		LastGeneralBalance:             o.LastGeneralBalance,
		LastGeneralNonce:               o.LastGeneralNonce,
		LastEscrowActiveBalance:        o.LastEscrowActiveBalance,
		LastEscrowActiveTotalShares:    o.LastEscrowActiveTotalShares,
		LastEscrowDebondingBalance:     o.LastEscrowDebondingBalance,
		LastEscrowDebondingTotalShares: o.LastEscrowDebondingTotalShares,
	}

	if !e.Valid() {
		return nil, errors.NewErrorFromMessage("account aggregator not valid", errors.NotValid)
	}

	return e, nil
}

func ToPersistence(ag *accountdomain.AccountAgg) (*orm.AccountAggModel, errors.ApplicationError) {
	if !ag.Valid() {
		return nil, errors.NewErrorFromMessage("account aggregator not valid", errors.NotValid)
	}

	return &orm.AccountAggModel{
		EntityModel: orm.EntityModel{
			ID: ag.ID,
		},
		AggregateModel: orm.AggregateModel{
			StartedAtHeight: ag.StartedAtHeight,
			StartedAt:       ag.StartedAt,
		},
		PublicKey:                      ag.PublicKey,
		LastGeneralBalance:             ag.LastGeneralBalance,
		LastGeneralNonce:               ag.LastGeneralNonce,
		LastEscrowActiveBalance:        ag.LastEscrowActiveBalance,
		LastEscrowActiveTotalShares:    ag.LastEscrowActiveTotalShares,
		LastEscrowDebondingBalance:     ag.LastEscrowDebondingBalance,
		LastEscrowDebondingTotalShares: ag.LastEscrowDebondingTotalShares,
	}, nil
}

func FromData(stateSyncable *syncabledomain.Syncable) ([]*accountdomain.AccountAgg, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var accounts []*accountdomain.AccountAgg
	for publicKey, info := range stateData.Data.Staking.Ledger {
		acc := &accountdomain.AccountAgg{
			DomainEntity: commons.NewDomainEntity(commons.EntityProps{}),
			Aggregate: commons.NewAggregate(commons.AggregateProps{
				StartedAtHeight: stateSyncable.Height,
				StartedAt:       stateSyncable.Time,
			}),

			PublicKey:                      types.PublicKey(publicKey.String()),
			LastGeneralBalance:             types.NewQuantity(info.General.Balance.ToBigInt()),
			LastGeneralNonce:               types.Nonce(info.General.Nonce),
			LastEscrowActiveBalance:        types.NewQuantity(info.Escrow.Active.Balance.ToBigInt()),
			LastEscrowActiveTotalShares:    types.NewQuantity(info.Escrow.Active.TotalShares.ToBigInt()),
			LastEscrowDebondingBalance:     types.NewQuantity(info.Escrow.Debonding.Balance.ToBigInt()),
			LastEscrowDebondingTotalShares: types.NewQuantity(info.Escrow.Debonding.TotalShares.ToBigInt()),
		}

		if !acc.Valid() {
			return nil, errors.NewErrorFromMessage("account aggregator not valid", errors.NotValid)
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}

