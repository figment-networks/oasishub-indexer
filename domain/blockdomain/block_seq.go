package blockdomain

import (
	"github.com/figment-networks/oasishub-indexer/domain/commons"
	"github.com/figment-networks/oasishub-indexer/domain/transactiondomain"
	"github.com/figment-networks/oasishub-indexer/domain/validatordomain"
	"github.com/figment-networks/oasishub-indexer/types"
)

type BlockSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	// Indexed data
	Hash              types.Hash      `json:"hash"`
	ProposerEntityUID types.PublicKey `json:"proposer_entity_uid"`
	AppVersion        int64           `json:"app_version"`
	BlockVersion      int64           `json:"block_version"`
	TransactionsCount types.Count     `json:"transactions_count"`

	// Associations
	Validators   []*validatordomain.ValidatorSeq     `json:"validators"`
	Transactions []*transactiondomain.TransactionSeq `json:"transactions"`
}

// - METHODS
func (b *BlockSeq) ValidOwn() bool {
	return b.AppVersion >= 0 &&
		b.BlockVersion >= 0
}

func (b *BlockSeq) EqualOwn(m BlockSeq) bool {
	return b.AppVersion == m.AppVersion &&
		b.BlockVersion == m.BlockVersion
}

func (b *BlockSeq) Valid() bool {
	return b.DomainEntity.Valid() &&
		b.Sequence.Valid() &&
		b.ValidOwn()
}

func (b *BlockSeq) Equal(m BlockSeq) bool {
	return b.DomainEntity.Equal(*m.DomainEntity) &&
		b.Sequence.Equal(*m.Sequence) &&
		b.EqualOwn(m)
}
