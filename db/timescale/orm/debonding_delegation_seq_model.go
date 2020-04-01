package orm

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type DebondingDelegationSeqModel struct {
	EntityModel
	SequenceModel

	ValidatorUID types.PublicKey
	DelegatorUID types.PublicKey
	Shares       types.Quantity
	DebondEnd    int64
}

func (DebondingDelegationSeqModel) TableName() string {
	return "debonding_delegation_sequences"
}
