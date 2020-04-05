package accountaggmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/accountagg"
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
)

func ToAggregate(stateSyncable *syncable.Model) ([]accountagg.Model, errors.ApplicationError) {
	stateData, err := syncablemapper.UnmarshalStateData(stateSyncable.Data)
	if err != nil {
		return nil, err
	}

	var accounts []accountagg.Model
	for publicKey, info := range stateData.Data.Staking.Ledger {
		acc := accountagg.Model{
			Aggregate: &shared.Aggregate{
				StartedAtHeight: stateSyncable.Height,
				StartedAt:       stateSyncable.Time,
			},

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

func ToView(s *accountagg.Model) map[string]interface{} {
	return map[string]interface{}{
		"id":                s.ID,
		"started_at_height": s.StartedAtHeight,
		"started_at":        s.StartedAt,

		"public_key":                         s.PublicKey,
		"last_general_balance":               s.LastGeneralBalance.String(),
		"last_general_nonce":                 s.LastGeneralNonce,
		"last_escrow_active_balance":         s.LastEscrowActiveBalance.String(),
		"last_escrow_active_total_shares":    s.LastEscrowActiveTotalShares.String(),
		"last_escrow_debonding_balance":      s.LastEscrowDebondingBalance.String(),
		"last_escrow_debonding_total_shares": s.LastEscrowDebondingTotalShares.String(),
	}
}
