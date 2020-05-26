package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type ValidatorSeq struct {
	*Model
	*Sequence

	EntityUID    string         `json:"entity_uid"`
	NodeUID      string         `json:"node_uid"`
	ConsensusUID string         `json:"consensus_uid"`
	Address      string         `json:"address"`
	Proposed     bool           `json:"proposed"`
	VotingPower  int64          `json:"voting_power"`
	TotalShares  types.Quantity `json:"total_shares"`
	// When precommit_validated is null it means that validator did not have chance to validate the block
	PrecommitValidated   *bool `json:"precommit_validated"`
	PrecommitBlockIDFlag int64 `json:"precommit_block_id_flag"`
	PrecommitIndex       int64 `json:"precommit_index"`
}

// - Methods
func (ValidatorSeq) TableName() string {
	return "validator_sequences"
}

func (vs *ValidatorSeq) Valid() bool {
	return vs.Sequence.Valid() &&
		vs.EntityUID != "" &&
		vs.NodeUID != "" &&
		vs.VotingPower >= 0 &&
		vs.TotalShares.Valid()
}

func (vs *ValidatorSeq) Equal(m ValidatorSeq) bool {
	return vs.Sequence.Equal(*m.Sequence) &&
		vs.EntityUID == m.EntityUID &&
		vs.NodeUID == m.NodeUID
}
