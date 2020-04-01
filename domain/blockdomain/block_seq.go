package blockdomain

import (
	"github.com/figment-networks/oasishub/domain/commons"
	"github.com/figment-networks/oasishub/types"
)

type BlockSeq struct {
	*commons.DomainEntity
	*commons.Sequence

	// Indexed data
	AppVersion        int64
	BlockVersion      int64
	TransactionsCount types.Count
}

// - METHODS
func (b *BlockSeq) ValidOwn() bool {
	return b.AppVersion >=0 &&
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



