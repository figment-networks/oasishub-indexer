package accountagg

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Aggregate

	PublicKey                         types.PublicKey `json:"public_key"`
	CurrentGeneralBalance             types.Quantity  `json:"current_general_balance"`
	CurrentGeneralNonce               types.Nonce     `json:"current_general_nonce"`
	CurrentEscrowActiveBalance        types.Quantity  `json:"current_escrow_active_balance"`
	CurrentEscrowActiveTotalShares    types.Quantity  `json:"current_escrow_active_total_shares"`
	CurrentEscrowDebondingBalance     types.Quantity  `json:"current_escrow_debonding_balance"`
	CurrentEscrowDebondingTotalShares types.Quantity  `json:"current_escrow_debonding_total_shares"`
}

// - METHODS
func (Model) TableName() string {
	return "account_aggregates"
}

func (aa *Model) ValidOwn() bool {
	return aa.PublicKey.Valid()
}

func (aa *Model) EqualOwn(m Model) bool {
	return aa.PublicKey.Equal(m.PublicKey)
}

func (aa *Model) Valid() bool {
	return aa.Model.Valid() &&
		aa.Aggregate.Valid() &&
		aa.ValidOwn()
}

func (aa *Model) Equal(m Model) bool {
	return aa.Model.Equal(*m.Model) &&
		aa.Aggregate.Equal(*m.Aggregate) &&
		aa.EqualOwn(m)
}

func (aa *Model) UpdateAggAttrs(u *Model) {
	aa.CurrentGeneralBalance = u.CurrentGeneralBalance
	aa.CurrentGeneralNonce = u.CurrentGeneralNonce
	aa.CurrentEscrowActiveBalance = u.CurrentEscrowActiveBalance
	aa.CurrentEscrowActiveBalance = u.CurrentEscrowActiveBalance
	aa.CurrentEscrowActiveTotalShares = u.CurrentEscrowActiveTotalShares
	aa.CurrentEscrowDebondingBalance = u.CurrentEscrowDebondingBalance
	aa.CurrentEscrowDebondingTotalShares = u.CurrentEscrowDebondingTotalShares
}
