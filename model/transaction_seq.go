package model

import (
	"github.com/figment-networks/oasishub-indexer/types"
)

type TransactionSeq struct {
	*Model
	*Sequence

	PublicKey string         `json:"public_key"`
	Hash      string         `json:"hash"`
	Nonce     uint64         `json:"nonce"`
	Fee       types.Quantity `json:"fee"`
	GasLimit  uint64         `json:"gas_limit"`
	GasPrice  types.Quantity `json:"gas_price"`
	Method    string         `json:"method"`
}

// - Methods
func (TransactionSeq) TableName() string {
	return "transaction_sequences"
}

func (ts *TransactionSeq) Valid() bool {
	return ts.Sequence.Valid() &&
		ts.PublicKey != "" &&
		ts.Hash != "" &&
		ts.Nonce >= 0
}

func (ts *TransactionSeq) Equal(m TransactionSeq) bool {
	return ts.Sequence.Equal(*m.Sequence) &&
		ts.PublicKey == m.PublicKey &&
		ts.Hash == m.Hash
}
