package model

type BlockSeq struct {
	*Model
	*Sequence

	// Indexed data
	Hash              string `json:"hash"`
	ProposerEntityUID string `json:"proposer_entity_uid"`
	AppVersion        int64  `json:"app_version"`
	BlockVersion      int64  `json:"block_version"`
	TransactionsCount int64  `json:"transactions_count"`
}

// - METHODS
func (BlockSeq) TableName() string {
	return "block_sequences"
}

func (b *BlockSeq) Valid() bool {
	return b.Sequence.Valid() &&
		b.AppVersion >= 0 &&
		b.BlockVersion >= 0
}

func (b *BlockSeq) Equal(m BlockSeq) bool {
	return b.Sequence.Equal(*m.Sequence) &&
		b.AppVersion == m.AppVersion &&
		b.BlockVersion == m.BlockVersion
}
