package validatordomain

import (
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/types"
)

type ValidatorSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	EntityUID    types.PublicKey
	NodeUID      types.PublicKey
	ConsensusUID types.PublicKey
	Address      string
	VotingPower  VotingPower
	Precommit    *Precommit
}

func (vs *ValidatorSeq) ValidOwn() bool {
	return vs.EntityUID.Valid() &&
		vs.NodeUID.Valid() &&
		vs.VotingPower.Valid()
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
