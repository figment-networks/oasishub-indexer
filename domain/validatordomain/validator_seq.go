package validatordomain

import (
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/types"
)

type ValidatorSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	EntityUID    types.PublicKey `json:"entity_uid"`
	NodeUID      types.PublicKey `json:"node_uid"`
	ConsensusUID types.PublicKey `json:"consensus_uid"`
	Address      string          `json:"address"`
	Proposed     bool            `json:"proposed"`
	VotingPower  VotingPower     `json:"voting_power"`
	TotalShares  types.Quantity  `json:"total_shares"`
	Precommit    *Precommit      `json:"precommit"`
}

func (vs *ValidatorSeq) ValidOwn() bool {
	return vs.EntityUID.Valid() &&
		vs.NodeUID.Valid() &&
		vs.VotingPower.Valid() &&
		vs.TotalShares.Valid()
}

func (vs *ValidatorSeq) EqualOwn(m ValidatorSeq) bool {
	return vs.EntityUID.Equal(m.EntityUID) &&
		vs.NodeUID.Equal(m.EntityUID) &&
		vs.Precommit.Equal(*m.Precommit)
}

func (vs *ValidatorSeq) Valid() bool {
	return vs.DomainEntity.Valid() &&
		vs.Sequence.Valid() &&
		vs.ValidOwn()
}

func (vs *ValidatorSeq) Equal(m ValidatorSeq) bool {
	return vs.DomainEntity.Equal(*m.DomainEntity) &&
		vs.Sequence.Equal(*m.Sequence) &&
		vs.EqualOwn(m)
}
