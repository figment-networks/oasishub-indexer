package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type ValidatorSeq struct {
	ID types.ID `json:"id"`

	*Sequence

	EntityUID    string         `json:"entity_uid"`
	Address      string         `json:"address"`
	Proposed     bool           `json:"proposed"`
	VotingPower  int64          `json:"voting_power"`
	TotalShares  types.Quantity `json:"total_shares"`
	// When precommit_validated is null it means that validator did not have chance to validate the block
	PrecommitValidated   *bool `json:"precommit_validated"`
}

// - Methods
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
