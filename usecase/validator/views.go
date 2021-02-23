package validator

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
)

type AggListView struct {
	Items []model.ValidatorAgg `json:"items"`
}

func ToAggListView(ms []model.ValidatorAgg) *AggListView {
	return &AggListView{
		Items: ms,
	}
}

type AggDetailsView struct {
	*model.Model
	*model.Aggregate

	Address                   string         `json:"address"`
	EntityUID                 string         `json:"entity_uid"`
	RecentTendermintAddress   string         `json:"recent_tendermint_address"`
	RecentVotingPower         int64          `json:"recent_voting_power"`
	RecentTotalShares         types.Quantity `json:"recent_total_shares"`
	RecentActiveEscrowBalance types.Quantity `json:"recent_active_escrow_balance"`
	RecentCommission          types.Quantity `json:"recent_commission"`
	RecentAsValidatorHeight   int64          `json:"recent_as_validator_height"`
	RecentProposedHeight      int64          `json:"recent_proposed_height"`
	AccumulatedProposedCount  int64          `json:"accumulated_proposed_count"`
	Uptime                    float64        `json:"uptime"`
	LogoURL                   string         `json:"logo_url"`
	EntityName                string         `json:"entity_name"`

	LastSequences []model.ValidatorSeq `json:"last_sequences"`
}

func ToAggDetailsView(m *model.ValidatorAgg, sequences []model.ValidatorSeq) *AggDetailsView {
	return &AggDetailsView{
		Model:     m.Model,
		Aggregate: m.Aggregate,

		Address:                   m.Address,
		EntityUID:                 m.EntityUID,
		RecentTendermintAddress:   m.RecentTendermintAddress,
		RecentVotingPower:         m.RecentVotingPower,
		RecentTotalShares:         m.RecentTotalShares,
		RecentActiveEscrowBalance: m.RecentActiveEscrowBalance,
		RecentAsValidatorHeight:   m.RecentAsValidatorHeight,
		RecentCommission:          m.RecentCommission,
		RecentProposedHeight:      m.RecentProposedHeight,
		AccumulatedProposedCount:  m.AccumulatedProposedCount,
		Uptime:                    float64(m.AccumulatedUptime) / float64(m.AccumulatedUptimeCount),
		LogoURL:                   m.LogoURL,
		EntityName:                m.EntityName,

		LastSequences: sequences,
	}
}

type SeqListItem struct {
	*model.Model
	*model.Sequence

	EntityUID           string         `json:"entity_uid"`
	Address             string         `json:"address"`
	Proposed            bool           `json:"proposed"`
	VotingPower         int64          `json:"voting_power"`
	TotalShares         types.Quantity `json:"total_shares"`
	ActiveEscrowBalance types.Quantity `json:"active_escrow_balance"`
	Commission          types.Quantity `json:"commission"`
	Rewards             types.Quantity `json:"rewards"`
	PrecommitValidated  *bool          `json:"precommit_validated"`
	EntityName          string         `json:"entity_name"`
}

type SeqListView struct {
	Items []SeqListItem `json:"items"`
}

func ToSeqListView(validatorSeqs []model.ValidatorSeq, validatorAggs []model.ValidatorAgg) SeqListView {
	nameLookup := make(map[string]string)
	for _, agg := range validatorAggs {
		nameLookup[agg.Address] = agg.EntityName
	}

	var items []SeqListItem
	for _, m := range validatorSeqs {

		item := SeqListItem{
			Sequence: m.Sequence,

			EntityUID:           m.EntityUID,
			Address:             m.Address,
			Proposed:            m.Proposed,
			VotingPower:         m.VotingPower,
			TotalShares:         m.TotalShares,
			ActiveEscrowBalance: m.ActiveEscrowBalance,
			Commission:          m.Commission,
			Rewards:             m.Rewards,
			PrecommitValidated:  m.PrecommitValidated,
		}

		if val, ok := nameLookup[m.Address]; ok {
			item.EntityName = val
		}

		items = append(items, item)
	}

	return SeqListView{
		Items: items,
	}
}
