package indexing

import (
	"context"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"reflect"
	"time"
)

func NewBlockFetcherTask(client *client.Client) pipeline.Task {
	return &BlockFetcherTask{
		client: client,
	}
}

type BlockFetcherTask struct {
	client *client.Client
}

func (t *BlockFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)
	block, err := t.client.Block.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.DebugJSON(block.GetBlock(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawBlock = block.GetBlock()
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

func (t *StateFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)
	state, err := t.client.State.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.DebugJSON(state.GetState(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawState = state.GetState()
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

func (t *ValidatorFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)
	validators, err := t.client.Validator.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

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

func (t *TransactionFetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)
	transactions, err := t.client.Transaction.GetByHeight(payload.CurrentHeight)
	if err != nil {
		return err
	}

	logger.DebugJSON(transactions.GetTransactions(),
		logger.Field("process", "pipeline"),
		logger.Field("stage", "fetcher"),
		logger.Field("request", "block"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.RawTransactions = transactions.GetTransactions()
	return nil
}
