package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type DelegationSeq struct {
	*Model
	*Sequence

	ValidatorUID string         `json:"validator_uid"`
	DelegatorUID string         `json:"delegator_uid"`
	Shares       types.Quantity `json:"shares"`
}

func (DelegationSeq) TableName() string {
	return "delegation_sequences"
}

func (d *DelegationSeq) Valid() bool {
	return d.Sequence.Valid() &&
		d.ValidatorUID != "" &&
		d.DelegatorUID != "" &&
		d.Shares.Valid()
}

func (d *DelegationSeq) Equal(m DelegationSeq) bool {
	return d.Sequence.Equal(*m.Sequence) &&
		d.ValidatorUID == m.ValidatorUID &&
		d.DelegatorUID == m.DelegatorUID
}
