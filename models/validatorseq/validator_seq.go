package validatorseq

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Sequence

	EntityUID    types.PublicKey `json:"entity_uid"`
	NodeUID      types.PublicKey `json:"node_uid"`
	ConsensusUID types.PublicKey `json:"consensus_uid"`
	Address      string          `json:"address"`
	Proposed     bool            `json:"proposed"`
	VotingPower  VotingPower     `json:"voting_power"`
	TotalShares  types.Quantity  `json:"total_shares"`
	// When precommit_validated is null it means that validator did not have chance to validate the block
	PrecommitValidated *bool  `json:"precommit_validated"`
	PrecommitType      *int64 `json:"precommit_type"`
	PrecommitIndex     *int64 `json:"precommit_index"`
}

// - Methods
func (Model) TableName() string {
	return "validator_sequences"
}

func (vs *Model) ValidOwn() bool {
	return vs.EntityUID.Valid() &&
		vs.NodeUID.Valid() &&
		vs.VotingPower.Valid() &&
		vs.TotalShares.Valid()
}

func (vs *Model) EqualOwn(m Model) bool {
	return vs.EntityUID.Equal(m.EntityUID) &&
		vs.NodeUID.Equal(m.EntityUID)
}

func (vs *Model) Valid() bool {
	return vs.Sequence.Valid() &&
		vs.ValidOwn()
}

func (vs *Model) Equal(m Model) bool {
	return vs.Model != nil &&
		m.Model != nil &&
		vs.Model.Equal(*m.Model) &&
		vs.Sequence.Equal(*m.Sequence) &&
		vs.EqualOwn(m)
}
