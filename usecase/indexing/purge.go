package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/indexer"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrPurgingDisabled = errors.New("purging disabled")
)

type purgeUseCase struct {
	cfg    *config.Config
	db     *store.Store
}

func NewPurgeUseCase(cfg *config.Config, db *store.Store) *purgeUseCase {
	return &purgeUseCase{
		cfg:    cfg,
		db:     db,
	}
}

func (uc *purgeUseCase) Execute(ctx context.Context) error {
	defer metric.LogUseCaseDuration(time.Now(), "purge")

	targetsReader, err := indexer.NewTargetsReader(uc.cfg.IndexerTargetsFile)
	if err != nil {
		return err
	}
	currentIndexVersion := targetsReader.GetCurrentVersion()

	if err := uc.purgeBlocks(currentIndexVersion); err != nil {
		return err
	}

	//if err := uc.purgeValidators(currentIndexVersion); err != nil {
	//	return err
	//}

	return nil
}


func (uc *purgeUseCase) purgeBlocks(currentIndexVersion int64) error {
	if err := uc.purgeBlockSequences(currentIndexVersion); uc.checkErr(err) {
		return err
	}
	//if err := uc.purgeBlockSummaries(types.IntervalHourly, uc.cfg.PurgeHourlySummariesInterval); uc.checkErr(err) {
	//	return err
	//}
	return nil
}

func (uc *purgeUseCase) purgeValidators(currentIndexVersion int64) error {
	if err := uc.purgeValidatorSequences(currentIndexVersion); uc.checkErr(err) {
		return err
	}
	if err := uc.purgeValidatorSummaries(types.IntervalHourly, uc.cfg.PurgeHourlySummariesInterval); uc.checkErr(err) {
		return err
	}
	return nil
}

func (uc *purgeUseCase) purgeBlockSequences(currentIndexVersion int64) error {
	blockSeq, err := uc.db.BlockSeq.FindMostRecent()
	if err != nil {
		return err
	}
	lastSeqTime := blockSeq.Time.Time

	duration, err := uc.parseDuration(uc.cfg.PurgeSequencesInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info("purging block sequences disabled. Purge interval set to 0.")
		}
		return err
	}

	purgeThresholdFromLastSeq := lastSeqTime.Add(- *duration)

	activityPeriods, err := uc.db.BlockSummary.FindActivityPeriods(types.IntervalDaily, currentIndexVersion)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("purging summarized block sequences... [older than=%s]", purgeThresholdFromLastSeq))

	deletedCount, err := uc.db.BlockSeq.DeleteOlderThan(purgeThresholdFromLastSeq, activityPeriods)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d block sequences purged", *deletedCount))

	return nil
}

func (uc *purgeUseCase) purgeBlockSummaries(interval types.SummaryInterval, purgeInterval string) error {
	blockSummary, err := uc.db.BlockSummary.FindMostRecentByInterval(interval)
	if err != nil {
		return err
	}
	lastSummaryTimeBucket := blockSummary.TimeBucket.Time

	duration, err := uc.parseDuration(purgeInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info(fmt.Sprintf("purging block summaries disabled [interval=%s] [purge_interval=%s]", interval, purgeInterval))
		}
		return err
	}

	purgeThreshold := lastSummaryTimeBucket.Add(- *duration)

	logger.Info(fmt.Sprintf("purging block summaries... [interval=%s] [older than=%s]", interval, purgeThreshold))

	deletedCount, err := uc.db.BlockSummary.DeleteOlderThan(interval, purgeThreshold)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d block summaries purged [interval=%s]", *deletedCount, interval))

	return nil
}

func (uc *purgeUseCase) purgeValidatorSequences(currentIndexVersion int64) error {
	validatorSeq, err := uc.db.ValidatorSeq.FindMostRecent()
	if err != nil {
		return err
	}
	lastSeqTime := validatorSeq.Time.Time

	blockSummary, err := uc.db.ValidatorSummary.FindMostRecent()
	if err != nil {
		return err
	}
	lastSummaryTimeBucket := blockSummary.TimeBucket.Time

	duration, err := uc.parseDuration(uc.cfg.PurgeSequencesInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info("purging validator sequences disabled. Purge interval set to 0.")
		}
		return err
	}

	purgeThresholdFromConfig := lastSeqTime.Add(- *duration)

	var purgeThreshold time.Time
	if purgeThresholdFromConfig.Before(lastSummaryTimeBucket) {
		purgeThreshold = purgeThresholdFromConfig
	} else {
		purgeThreshold = lastSummaryTimeBucket
	}

	logger.Info(fmt.Sprintf("purging validator sequences... [older than=%s]", purgeThreshold))

	deletedCount, err := uc.db.ValidatorSeq.DeleteOlderThan(purgeThreshold)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d validator sequences purged", *deletedCount))

	return nil
}

func (uc *purgeUseCase) purgeValidatorSummaries(interval types.SummaryInterval, purgeInterval string) error {
	blockSummary, err := uc.db.ValidatorSummary.FindMostRecentByInterval(interval)
	if err != nil {
		return err
	}
	lastSummaryTimeBucket := blockSummary.TimeBucket.Time

	duration, err := uc.parseDuration(purgeInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info(fmt.Sprintf("purging validator summaries disabled [interval=%s] [purge_interval=%s]", interval, purgeInterval))
		}
		return err
	}

	purgeThreshold := lastSummaryTimeBucket.Add(- *duration)

	logger.Info(fmt.Sprintf("purging validator summaries... [interval=%s] [older than=%s]", interval, purgeThreshold))

	deletedCount, err := uc.db.ValidatorSummary.DeleteOlderThan(interval, purgeThreshold)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d validator summaries purged [interval=%s]", *deletedCount, interval))

	return nil
}

func (uc *purgeUseCase) parseDuration(interval string) (*time.Duration, error) {
	duration, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}

	if duration == 0 {
		return nil, ErrPurgingDisabled
	}
	return &duration, nil
}

func (uc *purgeUseCase) checkErr(err error) bool {
	return err != nil && err != ErrPurgingDisabled && err != store.ErrNotFound
}