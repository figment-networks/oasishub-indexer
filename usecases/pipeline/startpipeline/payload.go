package startpipeline

import (
	"github.com/figment-networks/oasishub-indexer/models/accountagg"
	"github.com/figment-networks/oasishub-indexer/models/blockseq"
	"github.com/figment-networks/oasishub-indexer/models/debondingdelegationseq"
	"github.com/figment-networks/oasishub-indexer/models/delegationseq"
	"github.com/figment-networks/oasishub-indexer/models/stakingseq"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/models/transactionseq"
	"github.com/figment-networks/oasishub-indexer/models/validatoragg"
	"github.com/figment-networks/oasishub-indexer/models/validatorseq"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/pipeline"
	"sync"
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
	RetrievedAt   types.Time

	// Syncer stage
	BlockSyncable        *syncable.Model
	StateSyncable        *syncable.Model
	ValidatorsSyncable   *syncable.Model
	TransactionsSyncable *syncable.Model

	// Run stage
	CalculatedEntitiesData   map[string]CalculatedEntityData
	CalculatedValidatorsData map[string]CalculatedValidatorData
	CalculatedAccountsData   map[string]CalculatedAccountData

	// Aggregate stage
	NewAggregatedAccounts       []accountagg.Model
	UpdatedAggregatedAccounts   []accountagg.Model
	NewAggregatedValidators     []validatoragg.Model
	UpdatedAggregatedValidators []validatoragg.Model

	// Sequence stage
	BlockSequence                *blockseq.Model
	StakingSequence              *stakingseq.Model
	ValidatorSequences           []validatorseq.Model
	TransactionSequences         []transactionseq.Model
	DelegationSequences          []delegationseq.Model
	DebondingDelegationSequences []debondingdelegationseq.Model
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

	newP.NewAggregatedAccounts = append([]accountagg.Model(nil), p.NewAggregatedAccounts...)
	newP.UpdatedAggregatedAccounts = append([]accountagg.Model(nil), p.UpdatedAggregatedAccounts...)
	newP.NewAggregatedValidators = append([]validatoragg.Model(nil), p.NewAggregatedValidators...)
	newP.UpdatedAggregatedValidators = append([]validatoragg.Model(nil), p.UpdatedAggregatedValidators...)

	newP.BlockSequence = p.BlockSequence
	newP.StakingSequence = p.StakingSequence
	newP.ValidatorSequences = append([]validatorseq.Model(nil), p.ValidatorSequences...)
	newP.TransactionSequences = append([]transactionseq.Model(nil), p.TransactionSequences...)
	newP.DelegationSequences = append([]delegationseq.Model(nil), p.DelegationSequences...)
	newP.DebondingDelegationSequences = append([]debondingdelegationseq.Model(nil), p.DebondingDelegationSequences...)

	return newP
}

func (p *payload) MarkAsProcessed() {
	// Reset
	p.NewAggregatedAccounts = p.NewAggregatedAccounts[:0]
	p.UpdatedAggregatedAccounts = p.UpdatedAggregatedAccounts[:0]
	p.NewAggregatedValidators = p.NewAggregatedValidators[:0]
	p.UpdatedAggregatedValidators = p.UpdatedAggregatedValidators[:0]
	p.ValidatorSequences = p.ValidatorSequences[:0]
	p.TransactionSequences = p.TransactionSequences[:0]
	p.DelegationSequences = p.DelegationSequences[:0]
	p.DebondingDelegationSequences = p.DebondingDelegationSequences[:0]

	payloadPool.Put(p)
}
