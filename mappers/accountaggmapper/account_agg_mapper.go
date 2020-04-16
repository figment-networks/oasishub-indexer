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
	for publicKey, info := range stateData.Staking.Ledger {
		acc := accountagg.Model{
			Aggregate: &shared.Aggregate{
				StartedAtHeight: stateSyncable.Height,
				StartedAt:       stateSyncable.Time,
			},

			PublicKey:                         types.PublicKey(publicKey),
			CurrentGeneralBalance:             types.NewQuantityFromBytes(info.General.Balance),
			CurrentGeneralNonce:               types.Nonce(info.General.Nonce),
			CurrentEscrowActiveBalance:        types.NewQuantityFromBytes(info.Escrow.Active.Balance),
			CurrentEscrowActiveTotalShares:    types.NewQuantityFromBytes(info.Escrow.Active.TotalShares),
			CurrentEscrowDebondingBalance:     types.NewQuantityFromBytes(info.Escrow.Debonding.Balance),
			CurrentEscrowDebondingTotalShares: types.NewQuantityFromBytes(info.Escrow.Debonding.TotalShares),
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
