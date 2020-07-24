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

type indexingPipeline struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	targetsReader *targetsReader
	pipeline      *pipeline.Pipeline
}

func NewPipeline(cfg *config.Config, db *store.Store, client *client.Client) (*indexingPipeline, error) {
	p := pipeline.New(NewPayloadFactory())

	// Setup logger
	p.SetLogger(NewLogger())

	// Setup stage
	p.SetStage(
		pipeline.StageSetup,
		pipeline.SyncRunner(
			pipeline.RetryingTask(NewHeightMetaRetrieverTask(client), isTransient, 3),
		),
	)

	// Syncer stage
	p.SetStage(
		pipeline.StageSyncer,
		pipeline.SyncRunner(
			pipeline.RetryingTask(NewMainSyncerTask(db), isTransient, 3),
		),
	)

	// Fetcher stage
	p.SetStage(
		pipeline.StageFetcher,
		pipeline.AsyncRunner(
			pipeline.RetryingTask(NewBlockFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewStakingStateFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewStateFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewValidatorFetcherTask(client), isTransient, 3),
			pipeline.RetryingTask(NewTransactionFetcherTask(client), isTransient, 3),
		),
	)

	// Set parser stage
	p.SetStage(
		pipeline.StageParser,
		pipeline.AsyncRunner(
			NewBlockParserTask(),
			NewValidatorsParserTask(),
		),
	)

	// Set sequencer stage
	p.SetStage(
		pipeline.StageSequencer,
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
	p.SetStage(
		pipeline.StageAggregator,
		pipeline.AsyncRunner(
			pipeline.RetryingTask(NewAccountAggCreatorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewValidatorAggCreatorTask(db), isTransient, 3),
		),
	)

	// Add analyzer stage
	p.AddStageBefore(pipeline.StagePersistor, StageAnalyzer, pipeline.AsyncRunner(
		NewSystemEventCreatorTask(db.ValidatorSeq),
	), )

	// Set persistor stage
	p.SetStage(
		pipeline.StagePersistor,
		pipeline.AsyncRunner(
			pipeline.RetryingTask(NewSyncerPersistorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewBlockSeqPersistorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewValidatorSeqPersistorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewValidatorAggPersistorTask(db), isTransient, 3),
			pipeline.RetryingTask(NewSystemEventPersistorTask(db), isTransient, 3),
		),
	)

	// Create targets reader
	targetsReader, err := NewTargetsReader(cfg.IndexerTargetsFile)
	if err != nil {
		return nil, err
	}

	return &indexingPipeline{
		cfg:    cfg,
		db:     db,
		client: client,

		pipeline:      p,
		targetsReader: targetsReader,
	}, nil
}

type StartConfig struct {
	BatchSize   int64
	StartHeight int64
}

func (p *indexingPipeline) Start(ctx context.Context, startCfg StartConfig) error {
	currentIndexVersion := p.targetsReader.GetCurrentVersionID()

	source, err := NewIndexSource(p.cfg, p.db, p.client, &IndexSourceConfig{
		BatchSize:   startCfg.BatchSize,
		StartHeight: startCfg.StartHeight,
	})
	if err != nil {
		return err
	}

	sink := NewSink(p.db, currentIndexVersion)

	report, err := p.createReport(source.startHeight, source.endHeight, model.ReportKindIndex, currentIndexVersion)
	if err != nil {
		return err
	}

	ctxWithReport := context.WithValue(ctx, CtxReport, report)

	pipelineOptions, err := p.getPipelineOptions(pipelineOptionsConfig{})
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("starting pipeline [start=%d] [end=%d]", source.startHeight, source.endHeight))

	err = p.pipeline.Start(ctxWithReport, source, sink, pipelineOptions)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
	}

	logger.Info(fmt.Sprintf("pipeline completed [Err: %+v]", err))

	err = p.completeReport(report, source.Len(), sink.successCount, err)

	return err
}

type BackfillConfig struct {
	Parallel   bool
	Force      bool
	VersionIds []int64
	TargetIds  []int64
}

func (p *indexingPipeline) Backfill(ctx context.Context, backfillCfg BackfillConfig) error {
	currentIndexVersion := p.targetsReader.GetCurrentVersionID()

	source, err := NewBackfillSource(p.cfg, p.db, p.client, &BackfillSourceConfig{
		indexVersion: currentIndexVersion,
	})
	if err != nil {
		return err
	}

	sink := NewSink(p.db, currentIndexVersion)

	kind := model.ReportKindSequentialReindex
	if backfillCfg.Parallel {
		kind = model.ReportKindParallelReindex
	}

	if backfillCfg.Force {
		if err := p.db.Reports.DeleteReindexing(); err != nil {
			return err
		}
	}

	report, err := p.getReport(currentIndexVersion, source, kind)
	if err != nil {
		return err
	}

	if err := p.db.Syncables.SetProcessedAtForRange(report.ID, source.startHeight, source.endHeight); err != nil {
		return err
	}

	ctxWithReport := context.WithValue(ctx, CtxReport, report)

	pipelineOptions, err := p.getPipelineOptions(pipelineOptionsConfig{
		currentIndexVersion: currentIndexVersion,
		versionIds:          backfillCfg.VersionIds,
		targetIds:           backfillCfg.TargetIds,
	})
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("starting pipeline backfill [start=%d] [end=%d] [kind=%s]", source.startHeight, source.endHeight, kind))

	if err := p.pipeline.Start(ctxWithReport, source, sink, pipelineOptions); err != nil {
		return err
	}

	if err = p.completeReport(report, source.Len(), sink.successCount, err); err != nil {
		return err
	}

	logger.Info("pipeline backfill completed")

	return nil
}

type RunConfig struct {
	Height           int64
	DesiredVersionID int64
	DesiredTargetID  int64
	Dry              bool
}

func (p *indexingPipeline) Run(ctx context.Context, runCfg RunConfig) (*payload, error) {
	pipelineOptions, err := p.getPipelineOptions(pipelineOptionsConfig{
		dry: runCfg.Dry,
		versionIds: []int64{runCfg.DesiredVersionID},
		targetIds: []int64{runCfg.DesiredTargetID},
	})
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("running pipeline... [height=%d] [version=%d]", runCfg.Height, runCfg.DesiredTargetID))

	runPayload, err := p.pipeline.Run(ctx, runCfg.Height, pipelineOptions)
	if err != nil {
		metric.IndexerTotalErrors.Inc()
		logger.Info(fmt.Sprintf("pipeline completed with error [Err: %+v]", err))
		return nil, err
	}

	logger.Info("pipeline completed successfully")

	payload := runPayload.(*payload)
	return payload, nil
}

type pipelineOptionsConfig struct {
	dry                 bool
	currentIndexVersion int64
	versionIds          []int64
	targetIds           []int64
}

func (p *indexingPipeline) getPipelineOptions(optionsCfg pipelineOptionsConfig) (*pipeline.Options, error) {
	taskWhitelist, err := p.getTasksWhitelist(optionsCfg.currentIndexVersion, optionsCfg.versionIds, optionsCfg.targetIds)
	if err != nil {
		return nil, err
	}

	return &pipeline.Options{
		TaskWhitelist:   taskWhitelist,
		StagesBlacklist: p.getStagesBlacklist(optionsCfg.dry),
	}, nil
}

func (p *indexingPipeline) getTasksWhitelist(currentVersionId int64, versionIds []int64, targetIds []int64) ([]pipeline.TaskName, error) {
	var taskWhitelist []pipeline.TaskName

	if len(versionIds) == 0 && len(targetIds) == 0 {
		smallestIndexVersion, err := p.db.Syncables.FindSmallestIndexVersion()
		if err != nil {
			return nil, err
		}

		if currentVersionId == 0 || *smallestIndexVersion == currentVersionId {
			// No gaps in index versions
			// Do normal indexing including all available tasks
			taskWhitelist = p.targetsReader.GetAllTasks()
		} else {
			// There are records with smaller index versions
			// Include tasks from missing versions
			nextIndexVersion := *smallestIndexVersion + 1
			var ids []int64
			for i := nextIndexVersion; i <= currentVersionId; i++ {
				ids = append(ids, i)
			}

			taskWhitelist, err = p.targetsReader.GetTasksByVersionIds(ids)
			if err != nil {
				return nil, err
			}
		}

	} else {
		if len(versionIds) > 0 {
			tasks, err := p.targetsReader.GetTasksByVersionIds(versionIds)
			if err != nil {
				return nil, err
			}
			taskWhitelist = append(taskWhitelist, tasks...)
		}

		if len(targetIds) > 0 {
			tasks, err := p.targetsReader.GetTasksByTargetIds(targetIds)
			if err != nil {
				return nil, err
			}
			taskWhitelist = append(taskWhitelist, tasks...)
		}
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
