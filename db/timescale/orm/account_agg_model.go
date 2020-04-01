package orm

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type AccountAggModel struct {
	EntityModel
	AggregateModel

	// Associations
	//Account   AccountModel `gorm:"foreignkey"`
	//AccountID types.UUID

	PublicKey                      types.PublicKey
	LastGeneralBalance             types.Quantity
	LastGeneralNonce               types.Nonce
	LastEscrowActiveBalance        types.Quantity
	LastEscrowActiveTotalShares    types.Quantity
	LastEscrowDebondingBalance     types.Quantity
	LastEscrowDebondingTotalShares types.Quantity
}

func (AccountAggModel) TableName() string {
	return "account_aggregates"
}
