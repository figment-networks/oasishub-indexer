package validatoragg

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Aggregate

	EntityUID                types.PublicKey   `json:"entity_uid"`
	RecentAddress            string            `json:"recent_address"`
	RecentVotingPower        types.VotingPower `json:"recent_voting_power"`
	RecentTotalShares        types.Quantity    `json:"recent_total_shares"`
	RecentAsValidatorHeight  types.Height      `json:"recent_as_validator_height"`
	RecentProposedHeight     types.Height      `json:"recent_proposed_height"`
	AccumulatedProposedCount int64             `json:"accumulated_proposed_count"`
	AccumulatedUptime        int64             `json:"accumulated_uptime"`
	AccumulatedUptimeCount   int64             `json:"accumulated_uptime_count"`
}

// - METHODS
func (Model) TableName() string {
	return "validator_aggregates"
}

func (aa *Model) ValidOwn() bool {
	return aa.EntityUID.Valid()
}

func (aa *Model) EqualOwn(m Model) bool {
	return aa.EntityUID.Equal(m.EntityUID)
}

func (aa *Model) Valid() bool {
	return aa.Model.Valid() &&
		aa.Aggregate.Valid() &&
		aa.ValidOwn()
}

func (aa *Model) Equal(m Model) bool {
	return aa.Model != nil &&
		m.Model != nil &&
		aa.Model.Equal(*m.Model) &&
		aa.Aggregate.Equal(*m.Aggregate) &&
		aa.EqualOwn(m)
}

func (aa *Model) UpdateAggAttrs(entity Model) {
	aa.Aggregate.RecentAtHeight = entity.Aggregate.RecentAtHeight
	aa.Aggregate.RecentAt = entity.Aggregate.RecentAt

	aa.RecentVotingPower = entity.RecentVotingPower
	aa.RecentTotalShares = entity.RecentTotalShares
	aa.RecentAsValidatorHeight = entity.RecentAsValidatorHeight
	aa.RecentProposedHeight = entity.RecentProposedHeight
	aa.AccumulatedProposedCount = entity.AccumulatedProposedCount
	aa.AccumulatedUptimeCount = entity.AccumulatedUptimeCount
	aa.AccumulatedUptime = entity.AccumulatedUptime
}
