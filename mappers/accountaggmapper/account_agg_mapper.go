package accountaggmapper

import (
	"github.com/figment-networks/oasishub-indexer/mappers/syncablemapper"
	"github.com/figment-networks/oasishub-indexer/models/accountagg"
	"github.com/figment-networks/oasishub-indexer/models/debondingdelegationseq"
	"github.com/figment-networks/oasishub-indexer/models/delegationseq"
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

			PublicKey:                         types.PublicKey(publicKey.String()),
			CurrentGeneralBalance:             types.NewQuantity(info.General.Balance.ToBigInt()),
			CurrentGeneralNonce:               types.Nonce(info.General.Nonce),
			CurrentEscrowActiveBalance:        types.NewQuantity(info.Escrow.Active.Balance.ToBigInt()),
			CurrentEscrowActiveTotalShares:    types.NewQuantity(info.Escrow.Active.TotalShares.ToBigInt()),
			CurrentEscrowDebondingBalance:     types.NewQuantity(info.Escrow.Debonding.Balance.ToBigInt()),
			CurrentEscrowDebondingTotalShares: types.NewQuantity(info.Escrow.Debonding.TotalShares.ToBigInt()),
		}

		if !acc.Valid() {
			return nil, errors.NewErrorFromMessage("account aggregator not valid", errors.NotValid)
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}

type DetailsView struct {
	*accountagg.Model

	CurrentDelegations         []delegationseq.Model          `json:"current_delegations"`
	RecentDebondingDelegations []debondingdelegationseq.Model `json:"recent_debonding_delegations"`
}

func ToDetailsView(s *accountagg.Model, ds []delegationseq.Model, dds []debondingdelegationseq.Model) *DetailsView {
	return &DetailsView{
		Model:                      s,

		CurrentDelegations:         ds,
		RecentDebondingDelegations: dds,
	}
}
