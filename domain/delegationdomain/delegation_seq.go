package delegationdomain

import (
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/types"
)

type DelegationSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	ValidatorUID types.PublicKey
	DelegatorUID types.PublicKey
	Shares       types.Quantity
}

// - METHODS
func (d *DelegationSeq) ValidOwn() bool {
	return d.ValidatorUID.Valid() &&
		d.DelegatorUID.Valid() &&
		d.Shares.Valid()
}

func (d *DelegationSeq) EqualOwn(m DelegationSeq) bool {
	return d.ValidatorUID.Equal(m.ValidatorUID) &&
		d.DelegatorUID.Equal(m.DelegatorUID)
}

func (d *DelegationSeq) Valid() bool {
	return d.DomainEntity.Valid() &&
		d.Sequence.Valid() &&
		d.ValidOwn()
}

func (d *DelegationSeq) Equal(m DelegationSeq) bool {
	return d.ValidatorUID == m.ValidatorUID &&
		d.DelegatorUID == m.DelegatorUID &&
		d.EqualOwn(m)
}

