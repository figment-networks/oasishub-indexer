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
	"github.com/pkg/errors"
)

const (
	CtxReport = "context_report"

	StageAnalyzer = "AnalyzerStage"
)

var (
	ErrIndexCannotBeRun    = errors.New("cannot run index process")
	ErrBackfillCannotBeRun = errors.New("cannot run backfill process")
)

type indexingPipeline struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	targetsReader TargetsReader
	pipeline      *pipeline.Pipeline
}

func NewPipeline(cfg *config.Config, db *store.Store, client *client.Client) *indexingPipeline {
	return &indexingPipeline{
		cfg:    cfg,
		db:     db,
		client: client,
	}
}

type IndexCfg struct {
	StartHeight int64
	BatchSize   int64
}

func (p *indexingPipeline) Index(ctx context.Context, indexCfg IndexCfg) error {
	if err := p.canRunIndex(); err != nil {
		return err
	}

	currentIndexVersion, err := p.getCurrentIndexVersion()
	if err != nil {
		return err
	}

	source, err := NewIndexSource(p.cfg, p.db, p.client, indexCfg.StartHeight, indexCfg.BatchSize)
	if err != nil {
		return err
	}

	sink := NewSink(p.db, *currentIndexVersion)

	report, err := p.createReport(source.startHeight, source.endHeight, model.ReportKindIndex, *currentIndexVersion)
	if err != nil {
		return err
	}

	ctxWithReport := context.WithValue(ctx, CtxReport, report)

	// During indexing we always want to take tasks from all available versions
	pipelineOptions, err := p.getPipelineOptions(pipelineOptionsConfig{
		desiredVersionIds: p.targetsReader.GetAllVersionedVersionIds(),
	})
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("starting pipeline [start=%d] [end=%d] [options=%+v]", source.startHeight, source.endHeight, pipelineOptions))

	err = p.getPipeline().Start(ctxWithReport, source, sink, pipelineOptions)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
	}

	logger.Info(fmt.Sprintf("pipeline completed [Err: %+v]", err))

	err = p.completeReport(report, source.Len(), sink.successCount, err)

	return err
}

func (p *indexingPipeline) canRunIndex() error {
	isUpToDate, err := p.isUpToDate()
	if err != nil {
		return err
	}

	if !*isUpToDate {
		missingVersionIds, err := p.getMissingVersionIds()
		if err != nil {
			return err
		}

		if p.targetsReader.IsAnyVersionSequential(missingVersionIds) {
			return ErrIndexCannotBeRun
		}
	}
	return nil
}

type BackfillConfig struct {
	Parallel bool
	Force    bool
}

func (p *indexingPipeline) Backfill(ctx context.Context, backfillCfg BackfillConfig) error {
	if err := p.canRunBackfill(backfillCfg.Parallel); err != nil {
		return err
	}

	currentIndexVersion, err := p.getCurrentIndexVersion()
	if err != nil {
		return err
	}

	source, err := NewBackfillSource(p.cfg, p.db, p.client, *currentIndexVersion)
	if err != nil {
		return err
	}

	sink := NewSink(p.db, *currentIndexVersion)

	kind := model.ReportKindSequentialReindex
	if backfillCfg.Parallel {
		kind = model.ReportKindParallelReindex
	}

	if backfillCfg.Force {
		if err := p.db.Reports.DeleteReindexing(); err != nil {
			return err
		}
	}

	report, err := p.getReport(*currentIndexVersion, source, kind)
	if err != nil {
		return err
	}

	if err := p.db.Syncables.SetProcessedAtForRange(report.ID, source.startHeight, source.endHeight); err != nil {
		return err
	}

	ctxWithReport := context.WithValue(ctx, CtxReport, report)

	versionIds, err := p.getMissingVersionIds()
	if err != nil {
		return err
	}

	pipelineOptions, err := p.getPipelineOptions(pipelineOptionsConfig{
		desiredVersionIds: versionIds,
	})
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("starting pipeline backfill [start=%d] [end=%d] [kind=%s] [options=%+v]", source.startHeight, source.endHeight, kind, pipelineOptions))

	if err := p.getPipeline().Start(ctxWithReport, source, sink, pipelineOptions); err != nil {
		return err
	}

	if err = p.completeReport(report, source.Len(), sink.successCount, err); err != nil {
		return err
	}

	logger.Info("pipeline backfill completed")

	return nil
}

func (p *indexingPipeline) canRunBackfill(isParallel bool) error {
	isUpToDate, err := p.isUpToDate()
	if err != nil {
		return err
	}

	if !*isUpToDate {
		missingVersionIds, err := p.getMissingVersionIds()
		if err != nil {
			return err
		}

		if isParallel && p.targetsReader.IsAnyVersionSequential(missingVersionIds) {
			return ErrBackfillCannotBeRun
		}
	}

	return nil
}

func (p *indexingPipeline) isUpToDate() (*bool, error) {
	currentIndexVersion, err := p.getCurrentIndexVersion()
	if err != nil {
		return nil, err
	}

	var upToDate bool

	smallestIndexVersion, err := p.db.Syncables.FindSmallestIndexVersion()
	if err != nil {
		if err == store.ErrNotFound {
			// no records in database, we are definitely up to date
			upToDate = true
		} else {
			return nil, err
		}
	} else {
		if *smallestIndexVersion == *currentIndexVersion {
			upToDate = true
		}
	}

	return &upToDate, nil
}

func (p *indexingPipeline) getMissingVersionIds() ([]int64, error) {
	currentIndexVersion, err := p.getCurrentIndexVersion()
	if err != nil {
		return nil, err
	}

	var startIndexVersion int64

	smallestIndexVersion, err := p.db.Syncables.FindSmallestIndexVersion()
	if err != nil {
		if err == store.ErrNotFound {
			// When syncables not found in databases, set start version to first
			startIndexVersion = 1
		} else {
			return nil, err
		}
	} else {
		if *smallestIndexVersion < *currentIndexVersion {
			// There are records with smaller index versions
			// Include tasks from missing versions
			startIndexVersion = *smallestIndexVersion + 1
		} else if *smallestIndexVersion == *currentIndexVersion {
			// When everything is up to date, set start version to first
			startIndexVersion = 1
		} else {
			return nil, errors.New(fmt.Sprintf("current index version %d is too small", currentIndexVersion))
		}
	}

	var ids []int64
	for i := startIndexVersion; i <= *currentIndexVersion; i++ {
		ids = append(ids, i)
	}

	return ids, nil
}

type RunConfig struct {
	Height           int64
	DesiredVersionID int64
	DesiredTargetID  int64
	Dry              bool
}

func (p *indexingPipeline) Run(ctx context.Context, runCfg RunConfig) (*payload, error) {
	pipelineOptions, err := p.getPipelineOptions(pipelineOptionsConfig{
		dry:               runCfg.Dry,
		desiredVersionIds: []int64{runCfg.DesiredVersionID},
		desiredTargetIds:  []int64{runCfg.DesiredTargetID},
	})
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("running pipeline... [height=%d] [version=%d]", runCfg.Height, runCfg.DesiredTargetID))

	runPayload, err := p.getPipeline().Run(ctx, runCfg.Height, pipelineOptions)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
		logger.Info(fmt.Sprintf("pipeline completed with error [Err: %+v]", err))
		return nil, err
	}

	logger.Info("pipeline completed successfully")

	payload := runPayload.(*payload)
	return payload, nil
}

func (p *indexingPipeline) getPipeline() *pipeline.Pipeline {
	if p.pipeline == nil {
		defaultPipeline := pipeline.New(NewPayloadFactory())

		// Setup logger
		defaultPipeline.SetLogger(NewLogger())

		// Setup stage
		defaultPipeline.SetStage(
			pipeline.StageSetup,
			pipeline.SyncRunner(
				pipeline.RetryingTask(NewHeightMetaRetrieverTask(p.client), isTransient, 3),
			),
		)

		// Syncer stage
		defaultPipeline.SetStage(
			pipeline.StageSyncer,
			pipeline.SyncRunner(
				pipeline.RetryingTask(NewMainSyncerTask(p.db), isTransient, 3),
			),
		)

		// Fetcher stage
		defaultPipeline.SetStage(
			pipeline.StageFetcher,
			pipeline.AsyncRunner(
				pipeline.RetryingTask(NewBlockFetcherTask(p.client), isTransient, 3),
				pipeline.RetryingTask(NewStakingStateFetcherTask(p.client), isTransient, 3),
				pipeline.RetryingTask(NewStateFetcherTask(p.client), isTransient, 3),
				pipeline.RetryingTask(NewValidatorFetcherTask(p.client), isTransient, 3),
				pipeline.RetryingTask(NewTransactionFetcherTask(p.client), isTransient, 3),
			),
		)

		// Set parser stage
		defaultPipeline.SetStage(
			pipeline.StageParser,
			pipeline.AsyncRunner(
				NewBlockParserTask(),
				NewValidatorsParserTask(),
			),
		)

		// Set sequencer stage
		defaultPipeline.SetStage(
			pipeline.StageSequencer,
			pipeline.AsyncRunner(
				pipeline.RetryingTask(NewBlockSeqCreatorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewTransactionSeqCreatorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewStakingSeqCreatorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewValidatorSeqCreatorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewDelegationsSeqCreatorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewDebondingDelegationsSeqCreatorTask(p.db), isTransient, 3),
			),
		)

		// Set aggregator stage
		defaultPipeline.SetStage(
			pipeline.StageAggregator,
			pipeline.AsyncRunner(
				pipeline.RetryingTask(NewAccountAggCreatorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewValidatorAggCreatorTask(p.db), isTransient, 3),
			),
		)

		// Add analyzer stage
		defaultPipeline.AddStageBefore(pipeline.StagePersistor, StageAnalyzer, pipeline.AsyncRunner(
			NewSystemEventCreatorTask(p.cfg, p.db.ValidatorSeq),
		), )

		// Set persistor stage
		defaultPipeline.SetStage(
			pipeline.StagePersistor,
			pipeline.AsyncRunner(
				pipeline.RetryingTask(NewSyncerPersistorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewBlockSeqPersistorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewValidatorSeqPersistorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewValidatorAggPersistorTask(p.db), isTransient, 3),
				pipeline.RetryingTask(NewSystemEventPersistorTask(p.db), isTransient, 3),
			),
		)

		p.pipeline = defaultPipeline
	}
	return p.pipeline
}

func (p *indexingPipeline) getTargetsReader() (TargetsReader, error) {
	if p.targetsReader == nil {
		tr, err := NewTargetsReader(p.cfg.IndexerTargetsFile)
		if err != nil {
			return nil, err
		}
		p.targetsReader = tr
	}
	return p.targetsReader, nil
}

func (p *indexingPipeline) getCurrentIndexVersion() (*int64, error) {
	tr, err := p.getTargetsReader()
	if err != nil {
		return nil, err
	}
	currentVersionID := tr.GetCurrentVersionId()
	return &currentVersionID, nil
}

type pipelineOptionsConfig struct {
	desiredVersionIds []int64
	desiredTargetIds  []int64
	dry               bool
}

func (p *indexingPipeline) getPipelineOptions(optionsCfg pipelineOptionsConfig) (*pipeline.Options, error) {
	taskWhitelist, err := p.getTasksWhitelist(optionsCfg.desiredVersionIds, optionsCfg.desiredTargetIds)
	if err != nil {
		return nil, err
	}

	return &pipeline.Options{
		TaskWhitelist:   taskWhitelist,
		StagesBlacklist: p.getStagesBlacklist(optionsCfg.dry),
	}, nil
}

func (p *indexingPipeline) getTasksWhitelist(desiredVersionIds []int64, desiredTargetIds []int64) ([]pipeline.TaskName, error) {
	var taskWhitelist []pipeline.TaskName

	if len(desiredVersionIds) > 0 {
		tasks, err := p.targetsReader.GetTasksByVersionIds(desiredVersionIds)
		if err != nil {
			return nil, err
		}
		taskWhitelist = append(taskWhitelist, tasks...)
	}

	if len(desiredTargetIds) > 0 {
		tasks, err := p.targetsReader.GetTasksByTargetIds(desiredTargetIds)
		if err != nil {
			return nil, err
		}
		taskWhitelist = append(taskWhitelist, tasks...)
	}

	return getUniqueTaskNames(taskWhitelist), nil
}

func (p *indexingPipeline) getStagesBlacklist(dry bool) []pipeline.StageName {
	var stagesBlacklist []pipeline.StageName
	if dry {
		stagesBlacklist = append(stagesBlacklist, pipeline.StagePersistor)
	}
	return stagesBlacklist
}

func (p *indexingPipeline) getReport(indexVersion int64, source *backfillSource, kind model.ReportKind) (*model.Report, error) {
	report, err := p.db.Reports.FindNotCompletedByIndexVersion(indexVersion, model.ReportKindSequentialReindex, model.ReportKindParallelReindex)
	if err != nil {
		if err == store.ErrNotFound {
			report, err = p.createReport(source.startHeight, source.endHeight, kind, indexVersion)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if report.Kind != kind {
			return nil, errors.New(fmt.Sprintf("there is already reindexing in process [kind=%s] (use -force flag to override it)", report.Kind))
		}
	}
	return report, nil
}

func (p *indexingPipeline) createReport(startHeight int64, endHeight int64, kind model.ReportKind, indexVersion int64) (*model.Report, error) {
	report := &model.Report{
		Kind:         kind,
		IndexVersion: indexVersion,
		StartHeight:  startHeight,
		EndHeight:    endHeight,
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
