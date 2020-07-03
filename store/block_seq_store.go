package store

import (
	"fmt"
	"time"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/oasishub-indexer/model"
)

var (
	_ BlockSeqStore = (*blockSeqStore)(nil)
)

type BlockSeqStore interface {
	BaseStore

	FindBy(string, interface{}) (*model.BlockSeq, error)
	FindByHeight(int64) (*model.BlockSeq, error)
	GetAvgRecentTimes(int64) GetAvgRecentTimesResult
	FindMostRecent() (*model.BlockSeq, error)
	DeleteOlderThan(time.Time, []ActivityPeriodRow) (*int64, error)
	Summarize(types.SummaryInterval, []ActivityPeriodRow) ([]BlockSeqSummary, error)
}

func NewBlockSeqStore(db *gorm.DB) *blockSeqStore {
	return &blockSeqStore{scoped(db, model.BlockSeq{})}
}

// blockSeqStore handles operations on blocks
type blockSeqStore struct {
	baseStore
}

// FindBy returns a block for a matching attribute
func (s blockSeqStore) FindBy(key string, value interface{}) (*model.BlockSeq, error) {
	result := &model.BlockSeq{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByHeight returns a block with the matching height
func (s blockSeqStore) FindByHeight(height int64) (*model.BlockSeq, error) {
	return s.FindBy("height", height)
}

// GetAvgRecentTimesResult Contains results for GetAvgRecentTimes query
type GetAvgRecentTimesResult struct {
	StartHeight int64   `json:"start_height"`
	EndHeight   int64   `json:"end_height"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Count       int64   `json:"count"`
	Diff        float64 `json:"diff"`
	Avg         float64 `json:"avg"`
}

// GetAvgRecentTimes Gets average block times for recent blocks by limit
func (s *BlockSeqStore) GetAvgRecentTimes(limit int64) GetAvgRecentTimesResult {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels([]string{"BlockSeqStore_GetAvgRecentTimes"}))
	defer t.ObserveDuration()

	var res GetAvgRecentTimesResult
	s.db.Raw(blockTimesForRecentBlocksQuery, limit).Scan(&res)

	return res
}

// GetAvgTimesForIntervalRow Contains row of data for FindSummary query
type GetAvgTimesForIntervalRow struct {
	TimeInterval string  `json:"time_interval"`
	Count        int64   `json:"count"`
	Avg          float64 `json:"avg"`
}

// FindMostRecent finds most recent block sequence
func (s *blockSeqStore) FindMostRecent() (*model.BlockSeq, error) {
	blockSeq := &model.BlockSeq{}
	if err := findMostRecent(s.db, "time", blockSeq); err != nil {
		return nil, err
	}
	return blockSeq, nil
}

// DeleteOlderThan deletes block sequence older than given threshold
func (s *blockSeqStore) DeleteOlderThan(purgeThreshold time.Time, activityPeriods []ActivityPeriodRow) (*int64, error) {
	tx := s.db.
		Unscoped()

	hasIntervals := false
	for _, activityPeriod := range activityPeriods {
		// Make sure that there are many intervals (ie. days) in period
		if !activityPeriod.Min.Equal(activityPeriod.Max) {
			hasIntervals = true
			// Thus, we do not add 1 day to Max because we don't want to purge sequences within last day of period
			tx = tx.Where("time >= ? AND time < ?", activityPeriod.Min, activityPeriod.Max)
		}
	}

	if hasIntervals {
		tx.Where("time < ?", purgeThreshold).
			Delete(&model.BlockSeq{})

		if tx.Error != nil {
			return nil, checkErr(tx.Error)
		}
	} else {
		logger.Info("no block sequences to purge")
	}

	return &tx.RowsAffected, nil
}

type BlockSeqSummary struct {
	TimeBucket   types.Time `json:"time_bucket"`
	Count        int64      `json:"count"`
	BlockTimeAvg float64    `json:"block_time_avg"`
}

// Summarize gets the summarized version of block sequences
func (s *BlockSeqStore) Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]BlockSeqSummary, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels([]string{"BlockSummaryStore_Summarize"}))
	defer t.ObserveDuration()

	tx := s.db.
		Table(model.BlockSeq{}.TableName()).
		Select(summarizeBlocksQuerySelect, interval).
		Order("time_bucket").
		Group("time_bucket")

	if len(activityPeriods) == 1 {
		activityPeriod := activityPeriods[0]
		tx = tx.Or("time < ? OR time >= ?", activityPeriod.Min, activityPeriod.Max)
	} else {
		for i, activityPeriod := range activityPeriods {
			isLast := i == len(activityPeriods)-1

			if isLast {
				tx = tx.Or("time >= ?", activityPeriod.Max)
			} else {
				duration, err := time.ParseDuration(fmt.Sprintf("1%s", interval))
				if err != nil {
					return nil, err
				}
				tx = tx.Or("time >= ? AND time < ?", activityPeriod.Max.Add(duration), activityPeriods[i+1].Min)
			}
		}
	}

	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []BlockSeqSummary
	for rows.Next() {
		var summary BlockSeqSummary
		if err := s.db.ScanRows(rows, &summary); err != nil {
			return nil, err
		}

		models = append(models, summary)
	}
	return models, nil
}
