package indexer

import (
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/event/eventpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/model"
)

var (
	_ pipeline.PayloadFactory = (*payloadFactory)(nil)
	_ pipeline.Payload        = (*payload)(nil)
)

func NewPayloadFactory(constants *chainpb.GetConstantsResponse) *payloadFactory {
	return &payloadFactory{
		CommonPoolAddress: constants.GetCommonPoolAddress(),
	}
}

type payloadFactory struct {
	CommonPoolAddress string
}

func (pf *payloadFactory) GetPayload(currentHeight int64) pipeline.Payload {
	return &payload{
		CurrentHeight:     currentHeight,
		CommonPoolAddress: pf.CommonPoolAddress,
	}
}

type payload struct {
	CurrentHeight     int64
	CommonPoolAddress string

	// Setup stage
	HeightMeta HeightMeta

	// Syncer stage
	Syncable *model.Syncable

	// Fetcher stage
	RawBlock        *blockpb.Block
	RawEscrowEvents *eventpb.EscrowEvents
	RawStakingState *statepb.Staking
	RawState        *statepb.State
	RawTransactions []*transactionpb.Transaction
	RawValidators   []*validatorpb.Validator

	// Parser stage
	ParsedBlock      ParsedBlockData
	ParsedValidators ParsedValidatorsData
	BalanceEvents    []model.BalanceEvent

	// Aggregator stage
	NewAggregatedAccounts       []model.AccountAgg
	UpdatedAggregatedAccounts   []model.AccountAgg
	NewAggregatedValidators     []model.ValidatorAgg
	UpdatedAggregatedValidators []model.ValidatorAgg

	// Sequencer stage
	NewBlockSequence          *model.BlockSeq
	UpdatedBlockSequence      *model.BlockSeq
	NewValidatorSequences     []model.ValidatorSeq
	UpdatedValidatorSequences []model.ValidatorSeq

	StakingSequence              *model.StakingSeq
	TransactionSequences         []model.TransactionSeq
	DelegationSequences          []model.DelegationSeq
	DebondingDelegationSequences []model.DebondingDelegationSeq

	// Analyzer
	SystemEvents []*model.SystemEvent
}

func (p *payload) MarkAsProcessed() {}
