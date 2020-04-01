package orm

import (
	"github.com/figment-networks/oasishub/types"
)

type DelegationSeqModel struct {
	EntityModel
	SequenceModel

	ValidatorUID types.PublicKey
	DelegatorUID types.PublicKey
	Shares       types.Quantity
}

func (DelegationSeqModel) TableName() string {
	return "delegation_sequences"
}
