package model

import "github.com/figment-networks/oasishub-indexer/types"

type BlockSeq struct {
	ID types.ID `json:"id"`

	*Sequence

	// Indexed data
	TransactionsCount int64  `json:"transactions_count"`
}

// - METHODS
func (BlockSeq) TableName() string {
	return "block_sequences"
}

func (b *BlockSeq) Valid() bool {
	return b.Sequence.Valid() &&
		b.TransactionsCount >= 0
}

func (b *BlockSeq) Equal(m BlockSeq) bool {
	return b.Sequence.Equal(*m.Sequence) &&
		b.TransactionsCount == m.TransactionsCount
}

func (b *BlockSeq) Update(m BlockSeq) {
	b.TransactionsCount = m.TransactionsCount
}