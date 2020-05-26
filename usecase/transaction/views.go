package transaction

import (
	"github.com/figment-networks/oasishub-indexer/model"
)

type ListView struct {
	Items []model.TransactionSeq `json:"items"`
}

func ToListView(ts []model.TransactionSeq) *ListView {
	return &ListView{
		Items: ts,
	}
}
