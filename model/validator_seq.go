package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type ValidatorSeq struct {
	ID types.ID `json:"id"`

	*Sequence

	EntityUID          string         `json:"entity_uid"`
	Address            string         `json:"address"`
	Proposed           bool           `json:"proposed"`
	VotingPower        int64          `json:"voting_power"`
	TotalShares        types.Quantity `json:"total_shares"`
	PrecommitValidated *bool          `json:"precommit_validated"`
}

func (ValidatorSeq) TableName() string {
	return "validator_sequences"
}

func (vs *ValidatorSeq) Valid() bool {
	return vs.Sequence.Valid() &&
		vs.EntityUID != "" &&
		vs.VotingPower >= 0 &&
		vs.TotalShares.Valid()
}

func (vs *ValidatorSeq) Equal(m ValidatorSeq) bool {
	return vs.Sequence.Equal(*m.Sequence) &&
		vs.EntityUID == m.EntityUID
}

func (vs *ValidatorSeq) Update(m ValidatorSeq) {
	vs.EntityUID = m.EntityUID
	vs.Address = m.Address
	vs.Proposed = m.Proposed
	vs.VotingPower = m.VotingPower
	vs.TotalShares = m.TotalShares
	vs.PrecommitValidated = m.PrecommitValidated
}
