package indexing

import (
	"context"
	"fmt"
	"time"

	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

type summarizeUseCase struct {
	cfg *config.Config
	db  *store.Store
}

func NewSummarizeUseCase(cfg *config.Config, db *store.Store) *summarizeUseCase {
	return &summarizeUseCase{
		cfg: cfg,
		db:  db,
	}
}

func (uc *summarizeUseCase) Execute(ctx context.Context) error {
	defer metric.LogUseCaseDuration(time.Now(), "summarize")

	if err := uc.summarizeBlockSeq(types.IntervalHourly); err != nil {
		return err
	}

	if err := uc.summarizeBlockSeq(types.IntervalDaily); err != nil {
		return err
	}

	if err := uc.summarizeValidatorSeq(types.IntervalHourly); err != nil {
		return err
	}

	if err := uc.summarizeValidatorSeq(types.IntervalDaily); err != nil {
		return err
	}

	return nil
}

func (uc *summarizeUseCase) summarizeBlockSeq(interval types.SummaryInterval) error {
	logger.Info(fmt.Sprintf("summarizing block sequences... [interval=%s]", interval))

	last, err := uc.db.BlockSummary.FindMostRecentByInterval(interval)
	if err != nil {
		if err == store.ErrNotFound {
			last = nil
		}
	}

	rawSummaryItems, err := uc.db.BlockSeq.Summarize(interval, last)
	if err != nil {
		return err
	}

	var newModels []model.BlockSummary
	var existingModels []model.BlockSummary
	for _, rawSummary := range rawSummaryItems {
		summary := &model.Summary{
			TimeInterval: interval,
			TimeBucket:   rawSummary.TimeBucket,
		}
		query := model.BlockSummary{
			Summary: summary,
		}

		existingBlockSummary, err := uc.db.BlockSummary.Find(&query)
		if err != nil {
			if err == store.ErrNotFound {
				blockSummary := model.BlockSummary{
					Summary: summary,

					Count: rawSummary.Count,
					BlockTimeAvg:   rawSummary.BlockTimeAvg,
				}
				if err := uc.db.BlockSummary.Create(&blockSummary); err != nil {
					return err
				}
				newModels = append(newModels, blockSummary)
			} else {
				return err
			}
		} else {
			existingBlockSummary.Count = rawSummary.Count
			existingBlockSummary.BlockTimeAvg = rawSummary.BlockTimeAvg

			if err := uc.db.BlockSummary.Save(existingBlockSummary); err != nil {
				return err
			}
			existingModels = append(existingModels, *existingBlockSummary)
		}
	}

	logger.Info(fmt.Sprintf("block sequences summarized [created=%d] [updated=%d]", len(newModels), len(existingModels)))

	return nil
}

func (uc *summarizeUseCase) summarizeValidatorSeq(interval types.SummaryInterval) error {
	logger.Info(fmt.Sprintf("summarizing validator sequences... [interval=%s]", interval))

	last, err := uc.db.ValidatorSummary.FindMostRecentByInterval(interval)
	if err != nil {
		if err == store.ErrNotFound {
			last = nil
		}
	}

	rawSummaryItems, err := uc.db.ValidatorSeq.Summarize(interval, last)
	if err != nil {
		return err
	}

	var newModels []model.ValidatorSummary
	var existingModels []model.ValidatorSummary
	for _, rawSummary := range rawSummaryItems {
		summary := &model.Summary{
			TimeInterval: interval,
			TimeBucket:   rawSummary.TimeBucket,
		}
		query := model.ValidatorSummary{
			Summary: summary,

			EntityUID: rawSummary.EntityUID,
		}

		existingValidatorSummary, err := uc.db.ValidatorSummary.Find(&query)
		if err != nil {
			if err == store.ErrNotFound {
				validatorSummary := model.ValidatorSummary{
					Summary: summary,

					EntityUID:       rawSummary.EntityUID,
					VotingPowerAvg:  rawSummary.VotingPowerAvg,
					VotingPowerMax:  rawSummary.VotingPowerMax,
					VotingPowerMin:  rawSummary.VotingPowerMin,
					TotalSharesAvg:  rawSummary.TotalSharesAvg,
					TotalSharesMax:  rawSummary.TotalSharesMax,
					TotalSharesMin:  rawSummary.TotalSharesMin,
					ValidatedSum:    rawSummary.ValidatedSum,
					NotValidatedSum: rawSummary.NotValidatedSum,
					ProposedSum:     rawSummary.ProposedSum,
					UptimeAvg:       rawSummary.UptimeAvg,
				}

				if err := uc.db.ValidatorSummary.Create(&validatorSummary); err != nil {
					return err
				}
				newModels = append(newModels, validatorSummary)
			} else {
				return err
			}
		} else {
			existingValidatorSummary.VotingPowerAvg = rawSummary.VotingPowerAvg
			existingValidatorSummary.VotingPowerMax = rawSummary.VotingPowerMax
			existingValidatorSummary.VotingPowerMin = rawSummary.VotingPowerMin
			existingValidatorSummary.TotalSharesAvg = rawSummary.TotalSharesAvg
			existingValidatorSummary.TotalSharesMax = rawSummary.TotalSharesMax
			existingValidatorSummary.TotalSharesMin = rawSummary.TotalSharesMin
			existingValidatorSummary.ValidatedSum = rawSummary.ValidatedSum
			existingValidatorSummary.NotValidatedSum = rawSummary.NotValidatedSum
			existingValidatorSummary.ProposedSum = rawSummary.ProposedSum
			existingValidatorSummary.UptimeAvg = rawSummary.UptimeAvg

			if err := uc.db.ValidatorSummary.Save(existingValidatorSummary); err != nil {
				return err
			}
			existingModels = append(existingModels, *existingValidatorSummary)
		}
	}

	logger.Info(fmt.Sprintf("validator sequences summarized [created=%d] [updated=%d]", len(newModels), len(existingModels)))

	return nil
}
