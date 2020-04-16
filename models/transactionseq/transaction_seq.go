package transactionseq

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Sequence

	PublicKey types.PublicKey `json:"public_key"`
	Hash      types.Hash      `json:"hash"`
	Nonce     types.Nonce     `json:"nonce"`
	Fee       types.Quantity  `json:"fee"`
	GasLimit  uint64          `json:"gas_limit"`
	GasPrice  types.Quantity  `json:"gas_price"`
	Method    string          `json:"method"`
}

// - Methods
func (Model) TableName() string {
	return "transaction_sequences"
}

func (ts *Model) ValidOwn() bool {
	return ts.PublicKey.Valid() &&
		ts.Hash.Valid() &&
		ts.Nonce.Valid()
}

func (ts *Model) EqualOwn(m Model) bool {
	return ts.PublicKey.Equal(m.PublicKey) &&
		ts.Hash.Equal(m.Hash)
}

func (ts *Model) Valid() bool {
	return ts.Sequence.Valid() &&
		ts.ValidOwn()
}

func (ts *Model) Equal(m Model) bool {
	return ts.Model != nil &&
		m.Model != nil &&
		ts.Model.Equal(*m.Model) &&
		ts.Sequence.Equal(*m.Sequence) &&
		ts.EqualOwn(m)
}
