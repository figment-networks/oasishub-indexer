package accountagg

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Aggregate

	PublicKey                      types.PublicKey `json:"public_key"`
	LastGeneralBalance             types.Quantity  `json:"last_general_balance"`
	LastGeneralNonce               types.Nonce     `json:"last_general_nonce"`
	LastEscrowActiveBalance        types.Quantity  `json:"last_escrow_active_balance"`
	LastEscrowActiveTotalShares    types.Quantity  `json:"last_escrow_active_total_shares"`
	LastEscrowDebondingBalance     types.Quantity  `json:"last_escrow_debonding_balance"`
	LastEscrowDebondingTotalShares types.Quantity  `json:"last_escrow_debonding_total_shares"`
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
	aa.LastGeneralBalance = u.LastGeneralBalance
	aa.LastGeneralNonce = u.LastGeneralNonce
	aa.LastEscrowActiveBalance = u.LastEscrowActiveBalance
	aa.LastEscrowActiveBalance = u.LastEscrowActiveBalance
	aa.LastEscrowActiveTotalShares = u.LastEscrowActiveTotalShares
	aa.LastEscrowDebondingBalance = u.LastEscrowDebondingBalance
	aa.LastEscrowDebondingTotalShares = u.LastEscrowDebondingTotalShares
}
