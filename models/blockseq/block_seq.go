package blockseq

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Sequence

	// Indexed data
	Hash              types.Hash      `json:"hash"`
	ProposerEntityUID types.PublicKey `json:"proposer_entity_uid"`
	AppVersion        int64           `json:"app_version"`
	BlockVersion      int64           `json:"block_version"`
	TransactionsCount types.Count     `json:"transactions_count"`
}

// - METHODS
func (Model) TableName() string {
	return "block_sequences"
}

func (b *Model) ValidOwn() bool {
	return b.AppVersion >= 0 &&
		b.BlockVersion >= 0
}

func (b *Model) EqualOwn(m Model) bool {
	return b.AppVersion == m.AppVersion &&
		b.BlockVersion == m.BlockVersion
}

func (b *Model) Valid() bool {
	return b.Model.Valid() &&
		b.Sequence.Valid() &&
		b.ValidOwn()
}

func (b *Model) Equal(m Model) bool {
	return b.Model.Equal(*m.Model) &&
		b.Sequence.Equal(*m.Sequence) &&
		b.EqualOwn(m)
}
