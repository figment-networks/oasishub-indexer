package delegationdomain

import (
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/types"
)

type DebondingDelegationSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	ValidatorUID types.PublicKey `json:"validator_uid"`
	DelegatorUID types.PublicKey `json:"delegator_uid"`
	Shares       types.Quantity  `json:"shares"`
	DebondEnd    int64           `json:"debond_end"`
}

// - METHODS
func (d *DebondingDelegationSeq) ValidOwn() bool {
	return d.ValidatorUID.Valid() &&
		d.DelegatorUID.Valid() &&
		d.Shares.Valid()
}

func (d *DebondingDelegationSeq) Valid() bool {
	return d.DomainEntity.Valid() &&
		d.Sequence.Valid() &&
		d.ValidOwn()
}

func (d *DebondingDelegationSeq) EqualOwn(m DebondingDelegationSeq) bool {
	return d.ValidatorUID.Equal(m.ValidatorUID) &&
		d.DelegatorUID.Equal(m.DelegatorUID)
}

func (d *DebondingDelegationSeq) Equal(m DebondingDelegationSeq) bool {
	return d.ValidatorUID == m.ValidatorUID &&
		d.DelegatorUID == m.DelegatorUID &&
		d.EqualOwn(m)
}
