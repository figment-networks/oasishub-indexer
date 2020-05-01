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
			RecentGeneralBalance:             types.NewQuantityFromBytes(info.General.Balance),
			RecentGeneralNonce:               types.Nonce(info.General.Nonce),
			RecentEscrowActiveBalance:        types.NewQuantityFromBytes(info.Escrow.Active.Balance),
			RecentEscrowActiveTotalShares:    types.NewQuantityFromBytes(info.Escrow.Active.TotalShares),
			RecentEscrowDebondingBalance:     types.NewQuantityFromBytes(info.Escrow.Debonding.Balance),
			RecentEscrowDebondingTotalShares: types.NewQuantityFromBytes(info.Escrow.Debonding.TotalShares),
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
