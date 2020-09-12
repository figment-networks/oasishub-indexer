package model

import "github.com/figment-networks/oasishub-indexer/types"

type BalanceSummary struct {
	*Model
	*Summary

	Address       string           `json:"address"`
	EscrowAddress string           `json:"escrow_addr"`
	TotalAmount   types.Quantity   `json:"total_amount"`
	Kind          BalanceEventKind `json:"kind"`
}

type RawBalanceSummary struct {
	TimeBucket types.Time `json:"time_bucket"`
	// StartHeight   int64      `json:"min_height"`
	Kind          BalanceEventKind `json:"kind"`
	Address       string           `json:"address"`
	EscrowAddress string           `json:"escrow_address"`
	TotalAmount   types.Quantity   `json:"total_amount"`
}

func (BalanceSummary) TableName() string {
	return "balance_summary"
}

func (s *BalanceSummary) UpdateFromRaw(m RawBalanceSummary) {
	s.Address = m.Address
	s.EscrowAddress = m.EscrowAddress
	s.TotalAmount = m.TotalAmount
	s.Kind = m.Kind
}
