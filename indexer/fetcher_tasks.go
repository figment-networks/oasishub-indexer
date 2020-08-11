package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/metric"
	"time"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

const (
	BlockFetcherTaskName        = "BlockFetcher"
	RewardFetcherTaskName       = "RewardFetcher"
	StateFetcherTaskName        = "StateFetcher"
	StakingStateFetcherTaskName = "StakingStateFetcher"
	ValidatorFetcherTaskName    = "ValidatorFetcher"
	TransactionFetcherTaskName  = "TransactionFetcher"
)

func NewBlockFetcherTask(client *client.Client) pipeline.Task {
	return &BlockFetcherTask{
		client: client,
	}
}

type BlockFetcherTask struct {
	client *client.Client
}

func (t *BlockFetcherTask) GetName() string {
	return BlockFetcherTaskName
}

func (t *BlockFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	block, err := t.client.Block.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(block.GetBlock(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawBlock = block.GetBlock()
	return nil
}

type RewardFetcherTask struct {
	client *client.Client
}

func (t *RewardFetcherTask) GetName() string {
	return RewardFetcherTaskName
}

func (t *RewardFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	events, err := t.client.Event.GetRewardsByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(events.GetEvents(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "events"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawRewardEvents = events.GetEvents()
	return nil
}

func NewStateFetcherTask(client *client.Client) pipeline.Task {
	return &StateFetcherTask{
		client: client,
	}
}

type StateFetcherTask struct {
	client *client.Client
}

func (t *StateFetcherTask) GetName() string {
	return StateFetcherTaskName
}

func (t *StateFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)
	state, err := t.client.State.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(state.GetState(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawState = state.GetState()
	return nil
}

func NewStakingStateFetcherTask(client *client.Client) pipeline.Task {
	return &StakingStateFetcherTask{
		client: client,
	}
}

type StakingStateFetcherTask struct {
	client *client.Client
}

func (t *StakingStateFetcherTask) GetName() string {
	return StakingStateFetcherTaskName
}

func (t *StakingStateFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)
	state, err := t.client.State.GetStakingByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(state.GetStaking(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawStakingState = state.GetStaking()
	return nil
}

func NewValidatorFetcherTask(client *client.Client) pipeline.Task {
	return &ValidatorFetcherTask{
		client: client,
	}
}

type ValidatorFetcherTask struct {
	client *client.Client
}

func (t *ValidatorFetcherTask) GetName() string {
	return ValidatorFetcherTaskName
}

func (t *ValidatorFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)
	validators, err := t.client.Validator.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(validators.GetValidators(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawValidators = validators.GetValidators()
	return nil
}

func NewTransactionFetcherTask(client *client.Client) pipeline.Task {
	return &TransactionFetcherTask{
		client: client,
	}
}

type TransactionFetcherTask struct {
	client *client.Client
}

func (t *TransactionFetcherTask) GetName() string {
	return TransactionFetcherTaskName
}

func (t *TransactionFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)
	transactions, err := t.client.Transaction.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageFetcher, t.GetName(), payload.CurrentHeight))
	logger.DebugJSON(transactions.GetTransactions(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawTransactions = transactions.GetTransactions()
	return nil
}
