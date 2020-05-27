package block

import (
	"github.com/figment-networks/oasishub-indexer/model"
)

type DetailsView struct {
	*model.Model
	*model.Sequence

	Validators   []model.ValidatorSeq   `json:"validators"`
	Transactions []model.TransactionSeq `json:"transactions"`
}

func ToDetailsView(m *model.BlockSeq, vs []model.ValidatorSeq, ts []model.TransactionSeq) (*DetailsView, error) {
	return &DetailsView{
		Model:    m.Model,
		Sequence: m.Sequence,

		Validators:   vs,
		Transactions: ts,
	}, nil
}
