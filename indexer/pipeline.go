package indexer

import (
	"context"
	"fmt"
	"time"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

const (
	CtxReport = "context_report"

	StageAnalyzer = "AnalyzerStage"
)

var (
	ErrIsPristine          = errors.New("cannot run because database is empty")
	ErrIndexCannotBeRun    = errors.New("cannot run index process")
	ErrBackfillCannotBeRun = errors.New("cannot run backfill process")
)

var Now = time.Now

type indexingPipeline struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	pipeline     pipeline.DefaultPipeline
	status       *pipelineStatus
	configParser ConfigParser
}

func NewPipeline(cfg *config.Config, db *store.Store, client *client.Client) (*indexingPipeline, error) {
	constants, err := client.Chain.GetConstants()
	if err != nil {
		return nil, err
	}

	defaultPipeline := pipeline.NewDefault(NewPayloadFactory(constants))

	// Setup logger
	defaultPipeline.SetLogger(NewLogger())

	// Setup stage
	defaultPipeline.SetTasks(
		pipeline.StageSetup,
		pipeline.RetryingTask(NewHeightMetaRetrieverTask(client.Chain), isTransient, 3),
	)

	// Syncer stage
	defaultPipeline.SetTasks(
		pipeline.StageSyncer,
		pipeline.RetryingTask(NewMainSyncerTask(db.Syncables), isTransient, 3),
	)

	// Fetcher stage
	defaultPipeline.SetAsyncTasks(
		pipeline.StageFetcher,
		pipeline.RetryingTask(NewBlockFetcherTask(client.Block), isTransient, 3),
		pipeline.RetryingTask(NewStakingStateFetcherTask(client.State), isTransient, 3),
		pipeline.RetryingTask(NewStateFetcherTask(client.State), isTransient, 3),
		pipeline.RetryingTask(NewValidatorFetcherTask(client.Validator), isTransient, 3),
		pipeline.RetryingTask(NewTransactionFetcherTask(client.Transaction), isTransient, 3),
		pipeline.RetryingTask(NewEventsFetcherTask(client.Event), isTransient, 3),
	)

	// Set parser stage
	defaultPipeline.SetAsyncTasks(
		pipeline.StageParser,
		NewBlockParserTask(),
		NewValidatorsParserTask(),
	)

	// Set sequencer stage
	defaultPipeline.SetAsyncTasks(
		pipeline.StageSequencer,
		pipeline.RetryingTask(NewBlockSeqCreatorTask(db.BlockSeq), isTransient, 3),
		pipeline.RetryingTask(NewTransactionSeqCreatorTask(db.TransactionSeq), isTransient, 3),
		pipeline.RetryingTask(NewStakingSeqCreatorTask(db.StakingSeq), isTransient, 3),
		pipeline.RetryingTask(NewValidatorSeqCreatorTask(db.ValidatorSeq), isTransient, 3),
		pipeline.RetryingTask(NewDelegationsSeqCreatorTask(db.DelegationSeq), isTransient, 3),
		pipeline.RetryingTask(NewDebondingDelegationsSeqCreatorTask(db.DebondingDelegationSeq), isTransient, 3),
	)

	// Set aggregator stage
	defaultPipeline.SetAsyncTasks(
		pipeline.StageAggregator,
		pipeline.RetryingTask(NewAccountAggCreatorTask(db.AccountAgg), isTransient, 3),
		pipeline.RetryingTask(NewValidatorAggCreatorTask(db.ValidatorAgg), isTransient, 3),
	)

	// Add analyzer stage
	defaultPipeline.AddStageBefore(pipeline.StagePersistor, pipeline.NewStageWithTasks(StageAnalyzer, NewSystemEventCreatorTask(cfg, db.ValidatorSeq)))

	// Set persistor stage
	defaultPipeline.SetAsyncTasks(
		pipeline.StagePersistor,
		pipeline.RetryingTask(NewSyncerPersistorTask(db.Syncables), isTransient, 3),
		pipeline.RetryingTask(NewBlockSeqPersistorTask(db.BlockSeq), isTransient, 3),
		pipeline.RetryingTask(NewValidatorSeqPersistorTask(db.ValidatorSeq), isTransient, 3),
		pipeline.RetryingTask(NewValidatorAggPersistorTask(db.ValidatorAgg), isTransient, 3),
		pipeline.RetryingTask(NewSystemEventPersistorTask(db.SystemEvents), isTransient, 3),
	)

	configParser, err := NewConfigParser(cfg.IndexerConfigFile)
	if err != nil {
		return nil, err
	}

	statusChecker := pipelineStatusChecker{db.Syncables, configParser.GetCurrentVersionId()}
	pipelineStatus, err := statusChecker.getStatus()
	if err != nil {
		return nil, err
	}

	return &indexingPipeline{
		cfg:    cfg,
		db:     db,
		client: client,

		pipeline:     defaultPipeline,
		status:       pipelineStatus,
		configParser: configParser,
	}, nil
}

type IndexConfig struct {
	StartHeight int64
	BatchSize   int64
}

// Index starts indexing process
func (o *indexingPipeline) Index(ctx context.Context, indexCfg IndexConfig) error {
	if err := o.canRunIndex(); err != nil {
		return err
	}

	currentIndexVersion := o.configParser.GetCurrentVersionId()

	source, err := NewIndexSource(o.cfg, o.db.Syncables, o.client.Chain, indexCfg.StartHeight, indexCfg.BatchSize)
	if err != nil {
		return err
	}

	sink := NewSink(o.db, currentIndexVersion)

	reportCreator := &reportCreator{
		kind:         model.ReportKindIndex,
		indexVersion: currentIndexVersion,
		startHeight:  source.startHeight,
		endHeight:    source.endHeight,
		store:        o.db.Reports,
	}

	versionIds := o.configParser.GetAllVersionedVersionIds()
	pipelineOptionsCreator := &pipelineOptionsCreator{
		configParser: o.configParser,

		desiredVersionIds: versionIds,
	}
	pipelineOptions, err := pipelineOptionsCreator.parse()
	if err != nil {
		return err
	}

	if err := reportCreator.create(); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("starting pipeline [start=%d] [end=%d] [options=%+v]", source.startHeight, source.endHeight, pipelineOptions))

	ctxWithReport := context.WithValue(ctx, CtxReport, reportCreator.report)
	err = o.pipeline.Start(ctxWithReport, source, sink, pipelineOptions)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
	}

	logger.Info(fmt.Sprintf("pipeline completed [Err: %+v]", err))

	err = reportCreator.complete(source.Len(), sink.successCount, err)

	return nil
}

func (o *indexingPipeline) canRunIndex() error {
	if !o.status.isPristine && !o.status.isUpToDate {
		if o.configParser.IsAnyVersionSequential(o.status.missingVersionIds) {
			return ErrIndexCannotBeRun
		}
	}
	return nil
}

type BackfillConfig struct {
	Parallel bool
	Force    bool
}

// Backfill starts backfill process
func (o *indexingPipeline) Backfill(ctx context.Context, backfillCfg BackfillConfig) error {
	if err := o.canRunBackfill(backfillCfg.Parallel); err != nil {
		return err
	}

	currentIndexVersion := o.configParser.GetCurrentVersionId()
	kind := model.ReportKindSequentialReindex

	source, err := NewBackfillSource(o.cfg, o.db.Syncables, currentIndexVersion)
	if err != nil {
		return err
	}

	sink := NewSink(o.db, currentIndexVersion)

	if backfillCfg.Parallel {
		kind = model.ReportKindParallelReindex
	}
	if backfillCfg.Force {
		if err := o.db.Reports.DeleteByKinds([]model.ReportKind{model.ReportKindParallelReindex, model.ReportKindSequentialReindex}); err != nil {
			return err
		}
	}

	reportCreator := &reportCreator{
		kind:         kind,
		indexVersion: currentIndexVersion,
		startHeight:  source.startHeight,
		endHeight:    source.endHeight,
		store:        o.db.Reports,
	}

	versionIds := o.status.missingVersionIds
	pipelineOptionsCreator := &pipelineOptionsCreator{
		configParser:      o.configParser,
		desiredVersionIds: versionIds,
	}
	pipelineOptions, err := pipelineOptionsCreator.parse()
	if err != nil {
		return err
	}

	if err := o.db.Syncables.ResetProcessedAtForRange(source.startHeight, source.endHeight); err != nil {
		return err
	}

	if err := reportCreator.createIfNotExists(model.ReportKindSequentialReindex, model.ReportKindParallelReindex); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("starting pipeline [start=%d] [end=%d] [options=%+v]", source.startHeight, source.endHeight, pipelineOptions))

	ctxWithReport := context.WithValue(ctx, CtxReport, reportCreator.report)
	err = o.pipeline.Start(ctxWithReport, source, sink, pipelineOptions)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
	}

	logger.Info(fmt.Sprintf("pipeline completed [Err: %+v]", err))

	err = reportCreator.complete(source.Len(), sink.successCount, err)

	return nil
}

func (o *indexingPipeline) canRunBackfill(isParallel bool) error {
	if o.status.isPristine {
		return ErrIsPristine
	}

	if !o.status.isUpToDate {
		if isParallel && o.configParser.IsAnyVersionSequential(o.status.missingVersionIds) {
			return ErrBackfillCannotBeRun
		}
	}
	return nil
}

type RunConfig struct {
	Height            int64
	DesiredVersionIDs []int64
	DesiredTargetIDs  []int64
	Dry               bool
}

// Run runs pipeline just for one height
func (o *indexingPipeline) Run(ctx context.Context, runCfg RunConfig) (*payload, error) {
	pipelineOptionsCreator := &pipelineOptionsCreator{
		configParser:      o.configParser,
		dry:               runCfg.Dry,
		desiredVersionIds: runCfg.DesiredVersionIDs,
		desiredTargetIds:  runCfg.DesiredTargetIDs,
	}

	pipelineOptions, err := pipelineOptionsCreator.parse()
	if err != nil {
		return nil, err
	}

	logger.Info("running pipeline...",
		logger.Field("height", runCfg.Height),
		logger.Field("versions", runCfg.DesiredVersionIDs),
		logger.Field("targets", runCfg.DesiredTargetIDs),
	)

	runPayload, err := o.pipeline.Run(ctx, runCfg.Height, pipelineOptions)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
		logger.Info(fmt.Sprintf("pipeline completed with error [Err: %+v]", err))
		return nil, err
	}

	logger.Info("pipeline completed successfully")

	payload := runPayload.(*payload)
	return payload, nil
}

func isTransient(error) bool {
	return true
}
