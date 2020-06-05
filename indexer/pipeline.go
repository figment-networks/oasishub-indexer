package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"time"
)

const (
	CtxReport = "context_report"
)

type indexingPipeline struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	pipeline *pipeline.Pipeline
}

func NewPipeline(cfg *config.Config, db *store.Store, client *client.Client) (*indexingPipeline, error) {
	p := pipeline.New(NewPayloadFactory())

	// Set options to control what stages and what indexing tasks to execute
	p.SetOptions(&pipeline.Options{
		//StagesWhitelist: []pipeline.StageName{pipeline.StageFetcher},
		IndexingTasksBlacklist: []string{
			"accountAggCreatorTask",
			"transactionSeqCreatorTask",
			"stakingSeqCreatorTask",
			"delegationsSeqCreatorTask",
			"debondingDelegationsSeqCreatorTask",
		},
	})

	// Setup stage
	p.SetSetupStage(
		pipeline.SyncRunner(
			pipeline.RetryingTask(NewHeightMetaRetrieverTask(client), isTransient, 3),
		),
	)

	// Syncer stage
	p.SetSyncerStage(
		pipeline.SyncRunner(
			pipeline.RetryingTask(NewMainSyncerTask(db), isTransient, 3),
		),
	)

	// Fetcher stage
	p.SetFetcherStage(
		pipeline.AsyncRunner(
			pipeline.RetryingTask(NewBlockFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewStateFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewValidatorFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewTransactionFetcherTask(client), isTransient, 3),
		),
	)

	// Set parser stage
	p.SetParserStage(
		pipeline.AsyncRunner(
			NewParseBlockTask(),
			NewParseValidatorsTask(),
		),
	)

	// Set sequencer stage
	p.SetSequencerStage(
		pipeline.AsyncRunner(
			pipeline.RetryingTask(NewBlockSeqCreatorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewTransactionSeqCreatorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewStakingSeqCreatorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewValidatorSeqCreatorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewDelegationsSeqCreatorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewDebondingDelegationsSeqCreatorTask(db), isTransient, 3),
		),
	)

	// Set aggregator stage
	p.SetAggregatorStage(
		pipeline.AsyncRunner(
			//pipeline.RetryingTask(NewAccountAggCreatorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewValidatorAggCreatorTask(db), isTransient, 3),
		),
	)

	return &indexingPipeline{
		cfg:      cfg,
		db:       db,
		client:   client,
		pipeline: p,
	}, nil
}

func (p *indexingPipeline) Start(ctx context.Context, batchSize int64) error {
	source := NewSource(p.cfg, p.db, p.client, batchSize)
	sink := NewSink(p.db)

	logger.Info(fmt.Sprintf("Starting pipeline: %d - %d", source.startHeight, source.endHeight))

	report, err := p.createReport(source.startHeight, source.endHeight)
	if err != nil {
		return err
	}

	ctxWithReport := context.WithValue(ctx, CtxReport, report)

	err = p.pipeline.Start(ctxWithReport, source, sink)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
	}

	logger.Info(fmt.Sprintf("pipeline done [Err: %+v]", err))

	err = p.completeReport(report, source.Len(), sink.successCount, err)

	return err
}

func (p *indexingPipeline) createReport(startHeight int64, endHeight int64) (*model.Report, error) {
	report := &model.Report{
		StartHeight: startHeight,
		EndHeight:   endHeight,
	}
	if err := p.db.Reports.Create(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (p *indexingPipeline) completeReport(report *model.Report, totalCount int64, successCount int64, err error) error {
	report.Complete(successCount, totalCount - successCount, err)

	return p.db.Reports.Save(report)
}

func isTransient(error) bool {
	return true
}

func logTaskDuration(start time.Time, taskName string) {
	elapsed := time.Since(start)
	metric.IndexerTaskDuration.WithLabelValues(taskName).Set(elapsed.Seconds())
}
