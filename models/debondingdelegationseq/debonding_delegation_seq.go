package debondingdelegationseq

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Sequence

	ValidatorUID types.PublicKey `json:"validator_uid"`
	DelegatorUID types.PublicKey `json:"delegator_uid"`
	Shares       types.Quantity  `json:"shares"`
	DebondEnd    int64           `json:"debond_end"`
}

// - METHODS
func (Model) TableName() string {
	return "debonding_delegation_sequences"
}

func (d *Model) ValidOwn() bool {
	return d.ValidatorUID.Valid() &&
		d.DelegatorUID.Valid() &&
		d.Shares.Valid()
}

func (d *Model) Valid() bool {
	return d.Model.Valid() &&
		d.Sequence.Valid() &&
		d.ValidOwn()
}

func (d *Model) EqualOwn(m Model) bool {
	return d.ValidatorUID.Equal(m.ValidatorUID) &&
		d.DelegatorUID.Equal(m.DelegatorUID)
}

func (d *Model) Equal(m Model) bool {
	return d.Model.Equal(*m.Model) &&
		d.Sequence.Equal(*m.Sequence) &&
		d.EqualOwn(m)
}
