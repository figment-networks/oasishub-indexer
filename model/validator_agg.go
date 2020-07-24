package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type ValidatorAgg struct {
	*Model
	*Aggregate

	Address                  string         `json:"address"`
	EntityUID                string         `json:"entity_uid"`
	RecentTendermintAddress  string         `json:"recent_tendermint_address"`
	RecentVotingPower        int64          `json:"recent_voting_power"`
	RecentTotalShares        types.Quantity `json:"recent_total_shares"`
	RecentAsValidatorHeight  int64          `json:"recent_as_validator_height"`
	RecentProposedHeight     int64          `json:"recent_proposed_height"`
	AccumulatedProposedCount int64          `json:"accumulated_proposed_count"`
	AccumulatedUptime        int64          `json:"accumulated_uptime"`
	AccumulatedUptimeCount   int64          `json:"accumulated_uptime_count"`
	LogoURL                  string         `json:"logo_url"`
	EntityName               string         `json:"entity_name"`
}

func (ValidatorAgg) TableName() string {
	return "validator_aggregates"
}

func (aa *ValidatorAgg) Valid() bool {
	return aa.Aggregate.Valid() &&
		aa.EntityUID != ""
}

func (aa *ValidatorAgg) Equal(m ValidatorAgg) bool {
	return aa.Aggregate.Equal(*m.Aggregate) &&
		aa.EntityUID == m.EntityUID
}

func (aa *ValidatorAgg) Update(entity ValidatorAgg) {
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
