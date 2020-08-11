package model

import "github.com/figment-networks/oasishub-indexer/types"

type ValidatorSummary struct {
	*Model
	*Summary

	Address                string         `json:"address"`
	VotingPowerAvg         float64        `json:"voting_power_avg"`
	VotingPowerMax         float64        `json:"voting_power_max"`
	VotingPowerMin         float64        `json:"voting_power_min"`
	TotalSharesAvg         types.Quantity `json:"total_shares_avg"`
	TotalSharesMax         types.Quantity `json:"total_shares_max"`
	TotalSharesMin         types.Quantity `json:"total_shares_min"`
	ActiveEscrowBalanceAvg types.Quantity `json:"active_escrow_balance_avg"`
	ActiveEscrowBalanceMax types.Quantity `json:"active_escrow_balance_max"`
	ActiveEscrowBalanceMin types.Quantity `json:"active_escrow_balance_min"`
	CommissionAvg          types.Quantity `json:"commission_avg"`
	CommissionMax          types.Quantity `json:"commission_max"`
	CommissionMin          types.Quantity `json:"commission_min"`
	ValidatedSum           int64          `json:"validated_sum"`
	NotValidatedSum        int64          `json:"not_validated_sum"`
	ProposedSum            int64          `json:"proposed_sum"`
	UptimeAvg              float64        `json:"uptime_avg"`
}

func (ValidatorSummary) TableName() string {
	return "validator_summary"
}
