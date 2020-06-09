package indexer

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
			pipeline.RetryingTask(NewStakingStateFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewStateFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewValidatorFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewTransactionFetcherTask(client), isTransient, 3),
		),
	)

	// Set parser stage
	p.SetParserStage(
		pipeline.AsyncRunner(
			NewBlockParserTask(),
			NewValidatorsParserTask(),
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
			pipeline.RetryingTask(NewAccountAggCreatorTask(db), isTransient, 3),
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

type Options struct {
	BatchSize      int64
	Mode           pipeline.Mode
	CurrentVersion *int64
	DesiredVersion *int64
}

func (p *indexingPipeline) Start(ctx context.Context, indexingOptions Options) error {
	versionNumber, options, err := p.getIndexerOptions(indexingOptions)
	if err != nil {
		return err
	}

	source := NewSource(p.cfg, p.db, p.client, *versionNumber, indexingOptions.BatchSize)
	sink := NewSink(p.db, *versionNumber)

	logger.Info(fmt.Sprintf("starting pipeline [start=%d] [end=%d]", source.startHeight, source.endHeight))

	report, err := p.createReport(source.startHeight, source.endHeight)
	if err != nil {
		return err
	}

	ctxWithReport := context.WithValue(ctx, CtxReport, report)

	err = p.pipeline.Start(ctxWithReport, source, sink, options)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
	}

	logger.Info(fmt.Sprintf("pipeline completed [Err: %+v]", err))

	err = p.completeReport(report, source.Len(), sink.successCount, err)

	return err
}

func (p *indexingPipeline) getIndexerOptions(indexingOptions Options) (*int64, *pipeline.Options, error) {
	versionReader := pipeline.NewVersionReader(p.cfg.IndexerVersionsDir)

	//TODO: Use Up() and Version() based on mode for reindexing
	versionNumber, taskWhitelist, err := versionReader.All()
	if err != nil {
		return nil, nil, err
	}

	return versionNumber, &pipeline.Options{
		TaskWhitelist: taskWhitelist,
	}, nil
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
	report.Complete(successCount, totalCount-successCount, err)

	return p.db.Reports.Save(report)
}

func isTransient(error) bool {
	return true
}

