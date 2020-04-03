package orm

import (
	"github.com/figment-networks/oasishub-indexer/domain/validatordomain"
	"github.com/figment-networks/oasishub-indexer/types"
)

type ValidatorSeqModel struct {
	EntityModel
	SequenceModel

	// Associations
	//Validator   ValidatorModel `gorm:"foreignkey"`
	//ValidatorID types.UUID

	// Indexes
	EntityUID          types.PublicKey
	NodeUID            types.PublicKey
	ConsensusUID       types.PublicKey
	Address            string
	VotingPower        validatordomain.VotingPower
	TotalShares        types.Quantity
	Proposed           bool
	PrecommitValidated *bool
	PrecommitType      *int64
	PrecommitIndex     *int64
}

func (ValidatorSeqModel) TableName() string {
	return "validator_sequences"
}
