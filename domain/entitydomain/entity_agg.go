package entitydomain

import (
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/types"
)

type EntityAgg struct {
	*commons.DomainEntity
	*commons.Aggregate

	EntityUID types.PublicKey
}

// - METHODS
func (aa *EntityAgg) ValidOwn() bool {
	return aa.EntityUID.Valid()
}

func (aa *EntityAgg) EqualOwn(m EntityAgg) bool {
	return aa.EntityUID.Equal(m.EntityUID)
}

func (aa *EntityAgg) Valid() bool {
	return aa.DomainEntity.Valid() &&
		aa.Aggregate.Valid() &&
		aa.ValidOwn()
}

func (aa *EntityAgg) Equal(m EntityAgg) bool {
	return aa.DomainEntity.Equal(*m.DomainEntity) &&
		aa.Aggregate.Equal(*m.Aggregate) &&
		aa.EqualOwn(m)
}

func (aa *EntityAgg) Update(u *EntityAgg) {

}
