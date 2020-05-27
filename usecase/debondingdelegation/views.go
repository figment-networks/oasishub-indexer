package debondingdelegation

import (
	"github.com/figment-networks/oasishub-indexer/model"
)

type ListView struct {
	Items []model.DebondingDelegationSeq `json:"items"`
}

func ToListView(ms []model.DebondingDelegationSeq) *ListView {
	return &ListView{
		Items: ms,
	}
}
