package indexing

import (
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/model"
	"sync"
)

var (
	payloadPool = sync.Pool{
		New: func() interface{} {
			return new(payload)
		},
	}

	_ pipeline.PayloadFactory = (*payloadFactory)(nil)
	_ pipeline.Payload        = (*payload)(nil)
)

func NewPayloadFactory() *payloadFactory {
	return &payloadFactory{}
}

type payloadFactory struct{}

func (pf *payloadFactory) GetPayload() pipeline.Payload {
	return payloadPool.Get().(*payload)
}

type payload struct {
	CurrentHeight int64

	// Setup stage
	HeightMeta HeightMeta

	// Syncer stage
	Syncable *model.Syncable

	// Fetcher stage
	RawBlock        *blockpb.Block
	RawState        *statepb.State
	RawTransactions []*transactionpb.Transaction
	RawValidators   []*validatorpb.Validator

	// Parser stage
	ParsedBlock      ParsedBlockData
	ParsedValidators ParsedValidatorsData

	// Aggregator stage
	NewAggregatedAccounts       []model.AccountAgg
	UpdatedAggregatedAccounts   []model.AccountAgg
	NewAggregatedValidators     []model.ValidatorAgg
	UpdatedAggregatedValidators []model.ValidatorAgg

	// Sequencer stage
	BlockSequence                *model.BlockSeq
	StakingSequence              *model.StakingSeq
	ValidatorSequences           []model.ValidatorSeq
	TransactionSequences         []model.TransactionSeq
	DelegationSequences          []model.DelegationSeq
	DebondingDelegationSequences []model.DebondingDelegationSeq
}

func (p *payload) SetCurrentHeight(height int64) {
	p.CurrentHeight = height
}

func (p *payload) MarkAsProcessed() {
	payloadPool.Put(p)
}
