package store

import (
	"fmt"
	"time"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"
)

var (
	_ BlockSummaryStore = (*blockSummaryStore)(nil)
)

type BlockSummaryStore interface {
	BaseStore

	Find(*model.BlockSummary) (*model.BlockSummary, error)
	FindMostRecentByInterval(types.SummaryInterval) (*model.BlockSummary, error)
	FindActivityPeriods(types.SummaryInterval, int64) ([]ActivityPeriodRow, error)
	FindSummary(types.SummaryInterval, string) ([]model.BlockSummary, error)
	DeleteOlderThan(types.SummaryInterval, time.Time) (*int64, error)
}

func NewBlockSummaryStore(db *gorm.DB) *blockSummaryStore {
	return &blockSummaryStore{scoped(db, model.BlockSummary{})}
}

// blockSummaryStore handles operations on block summary
type blockSummaryStore struct {
	baseStore
}

// Find find block summary by query
func (s blockSummaryStore) Find(query *model.BlockSummary) (*model.BlockSummary, error) {
	var result model.BlockSummary

	err := s.db.
		Where(query).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindMostRecentByInterval finds most recent block summary for given time interval
func (s *blockSummaryStore) FindMostRecentByInterval(interval types.SummaryInterval) (*model.BlockSummary, error) {
	query := &model.BlockSummary{
		Summary: &model.Summary{TimeInterval: interval},
	}
	result := model.BlockSummary{}

	err := s.db.
		Where(query).
		Order("time_bucket DESC").
		Take(&result).
		Error

	return &result, checkErr(err)
}

type ActivityPeriodRow struct {
	Period int64
	Min    types.Time
	Max    types.Time
}

// FindActivityPeriods Finds activity periods
func (s *blockSummaryStore) FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]ActivityPeriodRow, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels("BlockSummaryStore_FindActivityPeriods"))
	defer t.ObserveDuration()

	query := getActivityPeriodsQuery(model.BlockSummary{}.TableName())

	var res []ActivityPeriodRow
	return res, s.db.Raw(query, fmt.Sprintf("1%s", interval), interval, indexVersion).Find(&res).Error
}

// FindSummary Gets summary of block sequences
func (s *blockSummaryStore) FindSummary(interval types.SummaryInterval, period string) ([]model.BlockSummary, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels("BlockSummaryStore_FindSummary"))
	defer t.ObserveDuration()

	var res []model.BlockSummary
	return res, s.db.Raw(allBlocksSummaryForIntervalQuery, interval, period, interval).Find(&res).Error
}

// DeleteOlderThan deletes block summary records older than given threshold
func (s *blockSummaryStore) DeleteOlderThan(interval types.SummaryInterval, purgeThreshold time.Time) (*int64, error) {
	res := s.db.
		Unscoped().
		Where("time_interval = ? AND time_bucket < ?", interval, purgeThreshold).
		Delete(&model.BlockSummary{})

	if res.Error != nil {
		return nil, checkErr(res.Error)
	}

	return &res.RowsAffected, nil
}
