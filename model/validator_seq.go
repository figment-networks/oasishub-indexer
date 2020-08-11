package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type ValidatorSeq struct {
	ID types.ID `json:"id"`

	*Sequence

	EntityUID           string         `json:"entity_uid"`
	Address             string         `json:"address"`
	Proposed            bool           `json:"proposed"`
	VotingPower         int64          `json:"voting_power"`
	TotalShares         types.Quantity `json:"total_shares"`
	ActiveEscrowBalance types.Quantity `json:"active_escrow_balance"`
	Commission          types.Quantity `json:"commission"`
	Rewards             types.Quantity `json:"rewards"`
	PrecommitValidated  *bool          `json:"precommit_validated"`
}

func (ValidatorSeq) TableName() string {
	return "validator_sequences"
}

func (vs *ValidatorSeq) Valid() bool {
	return vs.Sequence.Valid() &&
		vs.EntityUID != "" &&
		vs.VotingPower >= 0 &&
		vs.TotalShares.Valid() &&
		vs.ActiveEscrowBalance.Valid()
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
	vs.ActiveEscrowBalance = m.ActiveEscrowBalance
	vs.Commission = m.Commission
	vs.Rewards = m.Rewards
	vs.PrecommitValidated = m.PrecommitValidated
}
