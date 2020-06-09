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

	EntityUID                string         `json:"entity_uid"`
	RecentAddress            string         `json:"recent_address"`
	RecentVotingPower        int64          `json:"recent_voting_power"`
	RecentTotalShares        types.Quantity `json:"recent_total_shares"`
	RecentAsValidatorHeight  int64          `json:"recent_as_validator_height"`
	RecentProposedHeight     int64          `json:"recent_proposed_height"`
	AccumulatedProposedCount int64          `json:"accumulated_proposed_count"`
	Uptime                   float64        `json:"uptime"`

	LastSequences []model.ValidatorSeq `json:"last_sequences"`
}

func ToAggDetailsView(m *model.ValidatorAgg, sequences []model.ValidatorSeq) *AggDetailsView {
	return &AggDetailsView{
		Model:     m.Model,
		Aggregate: m.Aggregate,

		EntityUID:                m.EntityUID,
		RecentAddress:            m.RecentAddress,
		RecentVotingPower:        m.RecentVotingPower,
		RecentTotalShares:        m.RecentTotalShares,
		RecentAsValidatorHeight:  m.RecentAsValidatorHeight,
		RecentProposedHeight:     m.RecentProposedHeight,
		AccumulatedProposedCount: m.AccumulatedProposedCount,
		Uptime:                   float64(m.AccumulatedUptime) / float64(m.AccumulatedUptimeCount),

		LastSequences: sequences,
	}
}

type SeqListView struct {
	Items []model.ValidatorSeq `json:"items"`
}

func ToSeqListView(ms []model.ValidatorSeq) *SeqListView {
	return &SeqListView{
		Items: ms,
	}
}
