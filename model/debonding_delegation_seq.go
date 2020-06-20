package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type DebondingDelegationSeq struct {
	*Model
	*Sequence

	ValidatorUID string         `json:"validator_uid"`
	DelegatorUID string         `json:"delegator_uid"`
	Shares       types.Quantity `json:"shares"`
	DebondEnd    uint64         `json:"debond_end"`
}

func (DebondingDelegationSeq) TableName() string {
	return "debonding_delegation_sequences"
}

func (d *DebondingDelegationSeq) Valid() bool {
	return d.Sequence.Valid() &&
		d.ValidatorUID != "" &&
		d.DelegatorUID != "" &&
		d.Shares.Valid()
}

func (d *DebondingDelegationSeq) Equal(m DebondingDelegationSeq) bool {
	return d.Sequence.Equal(*m.Sequence) &&
		d.ValidatorUID == m.ValidatorUID &&
		d.DelegatorUID == m.DelegatorUID
}
