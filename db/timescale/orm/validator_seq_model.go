package orm

import (
	"github.com/figment-networks/oasishub/domain/validatordomain"
	"github.com/figment-networks/oasishub/types"
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
	PrecommitValidated *bool
	PrecommitType      *int64
	PrecommitIndex     *int64
}

func (ValidatorSeqModel) TableName() string {
	return "validator_sequences"
}
