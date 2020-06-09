package model

import "github.com/figment-networks/oasishub-indexer/types"

type ValidatorSummary struct {
	*Model
	*Summary

	EntityUID       string                `json:"entity_uid"`
	VotingPowerAvg  float64               `json:"voting_power_avg"`
	VotingPowerMax  float64               `json:"voting_power_max"`
	VotingPowerMin  float64               `json:"voting_power_min"`
	TotalSharesAvg  types.Quantity        `json:"total_shares_avg"`
	TotalSharesMax  types.Quantity        `json:"total_shares_max"`
	TotalSharesMin  types.Quantity        `json:"total_shares_min"`
	ValidatedSum    int64                 `json:"validated_sum"`
	NotValidatedSum int64                 `json:"not_validated_sum"`
	ProposedSum     int64                 `json:"proposed_sum"`
	UptimeAvg       float64               `json:"uptime_avg"`
}

func (ValidatorSummary) TableName() string {
	return "validator_summary"
}
