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

func (BalanceSummary) TableName() string {
	return "balance_summary"
}

func (s *BalanceSummary) Update(m BalanceSummary) {
	s.Address = m.Address
	s.EscrowAddress = m.EscrowAddress
	s.TotalAmount = m.TotalAmount
	s.Kind = m.Kind
}
