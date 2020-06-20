package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type AccountAgg struct {
	*Model
	*Aggregate

	PublicKey                        string `json:"public_key"`
	RecentGeneralBalance             types.Quantity  `json:"recent_general_balance"`
	RecentGeneralNonce               uint64     `json:"recent_general_nonce"`
	RecentEscrowActiveBalance        types.Quantity  `json:"recent_escrow_active_balance"`
	RecentEscrowActiveTotalShares    types.Quantity  `json:"recent_escrow_active_total_shares"`
	RecentEscrowDebondingBalance     types.Quantity  `json:"recent_escrow_debonding_balance"`
	RecentEscrowDebondingTotalShares types.Quantity  `json:"recent_escrow_debonding_total_shares"`
}

func (AccountAgg) TableName() string {
	return "account_aggregates"
}

func (aa *AccountAgg) Valid() bool {
	return aa.Aggregate.Valid() &&
		aa.PublicKey != ""
}

func (aa *AccountAgg) Equal(m AccountAgg) bool {
	return aa.Aggregate.Equal(*m.Aggregate) &&
		aa.PublicKey == m.PublicKey
}

func (aa *AccountAgg) Update(u *AccountAgg) {
	aa.Aggregate.RecentAtHeight = u.Aggregate.RecentAtHeight
	aa.Aggregate.RecentAt = u.Aggregate.RecentAt

	aa.RecentGeneralBalance = u.RecentGeneralBalance
	aa.RecentGeneralNonce = u.RecentGeneralNonce
	aa.RecentEscrowActiveBalance = u.RecentEscrowActiveBalance
	aa.RecentEscrowActiveTotalShares = u.RecentEscrowActiveTotalShares
	aa.RecentEscrowDebondingBalance = u.RecentEscrowDebondingBalance
	aa.RecentEscrowDebondingTotalShares = u.RecentEscrowDebondingTotalShares
}
