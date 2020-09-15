package model

import "github.com/figment-networks/oasishub-indexer/types"

type BalanceSummary struct {
	*Model
	*Summary

	StartHeight     int64          `json:"start_height"`
	Address         string         `json:"address"`
	EscrowAddress   string         `json:"escrow_address"`
	TotalRewards    types.Quantity `json:"total_rewards"`
	TotalCommission types.Quantity `json:"total_commission"`
	TotalSlashed    types.Quantity `json:"total_slashed"`
}

func (BalanceSummary) TableName() string {
	return "balance_summary"
}

func (s *BalanceSummary) Update(m BalanceSummary) {
	s.Address = m.Address
	s.EscrowAddress = m.EscrowAddress
	s.TotalRewards = m.TotalRewards
	s.TotalCommission = m.TotalCommission
	s.TotalSlashed = m.TotalSlashed
}
