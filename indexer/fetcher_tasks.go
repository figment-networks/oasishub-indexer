package indexer

import (
	"context"
	"fmt"
	"time"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

const (
	TaskNameBlockFetcher        = "BlockFetcher"
	TaskNameTransactionFetcher  = "TransactionFetcher"
	TaskNameEventFetcher        = "EventsFetcher"
	TaskNameStateFetcher        = "StateFetcher"
	TaskNameStakingStateFetcher = "StakingStateFetcher"
	TaskNameValidatorFetcher    = "ValidatorFetcher"
	TaskNameTransactionFetcher  = "TransactionFetcher"
)

func NewBlockFetcherTask(client client.BlockClient) pipeline.Task {
	return &BlockFetcherTask{
		client:         client,
		metricObserver: indexerTaskDuration.WithLabels([]string{TaskNameBlockFetcher}),
	}
}

type BlockFetcherTask struct {
	client         client.BlockClient
	metricObserver metrics.Observer
}

func (t *BlockFetcherTask) GetName() string {
	return TaskNameBlockFetcher
}

func (t *BlockFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)
	block, err := t.client.GetByHeight(payload.CurrentHeight)
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

func NewEventsFetcherTask(client client.EventClient) pipeline.Task {
	return &EventsFetcherTask{
		client: client,
	}
}

type EventsFetcherTask struct {
	client client.EventClient
}

func (t *EventsFetcherTask) GetName() string {
	return EventFetcherTaskName
}

func (t *EventsFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())
	payload := p.(*payload)

	events, err := t.client.GetAddEscrowEventsByHeight(payload.CurrentHeight)
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

	payload.RawEscrowEvents = events.GetEvents()
	return nil
}

func NewStateFetcherTask(client client.StateClient) pipeline.Task {
	return &StateFetcherTask{
		client:         client,
		metricObserver: indexerTaskDuration.WithLabels([]string{TaskNameStateFetcher}),
	}
}

type StateFetcherTask struct {
	client         client.StateClient
	metricObserver metrics.Observer
}

func (t *StateFetcherTask) GetName() string {
	return TaskNameStateFetcher
}

func (t *StateFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)
	state, err := t.client.GetByHeight(payload.CurrentHeight)
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

func NewStakingStateFetcherTask(client client.StateClient) pipeline.Task {
	return &StakingStateFetcherTask{
		client:         client,
		metricObserver: indexerTaskDuration.WithLabels([]string{TaskNameStakingStateFetcher}),
	}
}

type StakingStateFetcherTask struct {
	client         client.StateClient
	metricObserver metrics.Observer
}

func (t *StakingStateFetcherTask) GetName() string {
	return TaskNameStakingStateFetcher
}

func (t *StakingStateFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)
	state, err := t.client.GetStakingByHeight(payload.CurrentHeight)
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

func NewValidatorFetcherTask(client client.ValidatorClient) pipeline.Task {
	return &ValidatorFetcherTask{
		client:         client,
		metricObserver: indexerTaskDuration.WithLabels([]string{TaskNameValidatorFetcher}),
	}
}

type ValidatorFetcherTask struct {
	client         client.ValidatorClient
	metricObserver metrics.Observer
}

func (t *ValidatorFetcherTask) GetName() string {
	return TaskNameValidatorFetcher
}

func (t *ValidatorFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)
	validators, err := t.client.GetByHeight(payload.CurrentHeight)
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

func NewTransactionFetcherTask(client client.TransactionClient) pipeline.Task {
	return &TransactionFetcherTask{
		client:         client,
		metricObserver: indexerTaskDuration.WithLabels([]string{TaskNameTransactionFetcher}),
	}
}

type TransactionFetcherTask struct {
	client         client.TransactionClient
	metricObserver metrics.Observer
}

func (t *TransactionFetcherTask) GetName() string {
	return TaskNameTransactionFetcher
}

func (t *TransactionFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)
	transactions, err := t.client.GetByHeight(payload.CurrentHeight)
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
