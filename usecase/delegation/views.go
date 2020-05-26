package delegation

import (
	"github.com/figment-networks/oasishub-indexer/model"
)

type ListView struct {
	Items []model.DelegationSeq `json:"items"`
}

func ToListView(ms []model.DelegationSeq) *ListView {
	return &ListView{
		Items: ms,
	}
}
