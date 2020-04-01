package startpipeline

import (
	"github.com/figment-networks/oasishub/domain/accountdomain"
	"github.com/figment-networks/oasishub/domain/blockdomain"
	"github.com/figment-networks/oasishub/domain/delegationdomain"
	"github.com/figment-networks/oasishub/domain/entitydomain"
	"github.com/figment-networks/oasishub/domain/stakingdomain"
	"github.com/figment-networks/oasishub/domain/syncabledomain"
	"github.com/figment-networks/oasishub/domain/transactiondomain"
	"github.com/figment-networks/oasishub/domain/validatordomain"
	"github.com/figment-networks/oasishub/types"
	"github.com/figment-networks/oasishub/utils/pipeline"
	"sync"
	"time"
)

var (
	payloadPool = sync.Pool{
		New: func() interface{} {
			return new(payload)
		},
	}
)

type payload struct {
	StartHeight   types.Height
	EndHeight     types.Height
	CurrentHeight types.Height
	RetrievedAt   time.Time

	BlockSyncable        *syncabledomain.Syncable
	StateSyncable        *syncabledomain.Syncable
	ValidatorsSyncable   *syncabledomain.Syncable
	TransactionsSyncable *syncabledomain.Syncable

	NewAggregatedAccounts     []*accountdomain.AccountAgg
	UpdatedAggregatedAccounts []*accountdomain.AccountAgg
	NewAggregatedEntities     []*entitydomain.EntityAgg
	UpdatedAggregatedEntities []*entitydomain.EntityAgg

	BlockSequence                *blockdomain.BlockSeq
	ValidatorSequences           []*validatordomain.ValidatorSeq
	TransactionSequences         []*transactiondomain.TransactionSeq
	StakingSequence              *stakingdomain.StakingSeq
	DelegationSequences          []*delegationdomain.DelegationSeq
	DebondingDelegationSequences []*delegationdomain.DebondingDelegationSeq
}

func (p *payload) Clone() pipeline.Payload {
	newP := payloadPool.Get().(*payload)

	newP.StartHeight = p.StartHeight
	newP.EndHeight = p.EndHeight
	newP.RetrievedAt = p.RetrievedAt
	newP.BlockSyncable = p.BlockSyncable
	newP.StateSyncable = p.StateSyncable
	newP.ValidatorsSyncable = p.ValidatorsSyncable
	newP.TransactionsSyncable = p.TransactionsSyncable

	return newP
}

func (p *payload) MarkAsProcessed() {
	// Reset
	payloadPool.Put(p)
}
