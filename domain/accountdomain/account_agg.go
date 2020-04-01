package accountdomain

import (
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/types"
)

type AccountAgg struct {
	*commons.DomainEntity
	*commons.Aggregate

	PublicKey                      types.PublicKey
	LastGeneralBalance             types.Quantity
	LastGeneralNonce               types.Nonce
	LastEscrowActiveBalance        types.Quantity
	LastEscrowActiveTotalShares    types.Quantity
	LastEscrowDebondingBalance     types.Quantity
	LastEscrowDebondingTotalShares types.Quantity
}

// - METHODS
func (aa *AccountAgg) ValidOwn() bool {
	return aa.PublicKey.Valid()
}

func (aa *AccountAgg) EqualOwn(m AccountAgg) bool {
	return aa.PublicKey.Equal(m.PublicKey)
}

func (aa *AccountAgg) Valid() bool {
	return aa.DomainEntity.Valid() &&
		aa.Aggregate.Valid() &&
		aa.ValidOwn()
}

func (aa *AccountAgg) Equal(m AccountAgg) bool {
	return aa.DomainEntity.Equal(*m.DomainEntity) &&
		aa.Aggregate.Equal(*m.Aggregate) &&
		aa.EqualOwn(m)
}

func (aa *AccountAgg) Update(u *AccountAgg) {
	aa.LastGeneralBalance = u.LastGeneralBalance
	aa.LastGeneralNonce = u.LastGeneralNonce
	aa.LastEscrowActiveBalance = u.LastEscrowActiveBalance
	aa.LastEscrowActiveBalance = u.LastEscrowActiveBalance
	aa.LastEscrowActiveTotalShares = u.LastEscrowActiveTotalShares
	aa.LastEscrowDebondingBalance = u.LastEscrowDebondingBalance
	aa.LastEscrowDebondingTotalShares = u.LastEscrowDebondingTotalShares
}
