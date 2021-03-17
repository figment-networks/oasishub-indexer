package store

import (
	"fmt"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"
)

var (
	_ BalanceSummaryStore = (*balanceSummaryStore)(nil)
)

type BalanceSummaryStore interface {
	BaseStore

	Find(*model.BalanceSummary) (*model.BalanceSummary, error)
	GetSummariesByInterval(interval types.SummaryInterval, address string, start, end *types.Time) ([]model.BalanceSummary, error)
	FindActivityPeriods(types.SummaryInterval, int64) ([]ActivityPeriodRow, error)
}

func NewBalanceSummaryStore(db *gorm.DB) *balanceSummaryStore {
	return &balanceSummaryStore{scoped(db, model.BalanceSummary{})}
}

type balanceSummaryStore struct {
	baseStore
}

// Find find balance summary by query
func (s balanceSummaryStore) Find(query *model.BalanceSummary) (*model.BalanceSummary, error) {
	var result model.BalanceSummary

	err := s.db.
		Where(query).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindActivityPeriods Finds activity periods
func (s *balanceSummaryStore) FindActivityPeriods(interval types.SummaryInterval, indexVersion int64) ([]ActivityPeriodRow, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels("BalanceSummaryStore_FindActivityPeriods"))
	defer t.ObserveDuration()

	query := getActivityPeriodsQuery(model.BalanceSummary{}.TableName())

	var res []ActivityPeriodRow
	return res, s.db.Raw(query, fmt.Sprintf("1%s", interval), interval, indexVersion).Find(&res).Error
}

// GetSummariesByInterval Gets daily summary of balance events
func (s *balanceSummaryStore) GetSummariesByInterval(interval types.SummaryInterval, address string, start, end *types.Time) ([]model.BalanceSummary, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels("BalanceSummaryStore_GetDailySummaries"))
	defer t.ObserveDuration()

	tx := s.db.
		Table(model.BalanceSummary{}.TableName()).
		Select("*").
		Where("address = ? AND time_interval= ?", address, interval).
		Order("time_bucket")

	if !end.IsZero() {
		tx = tx.Where("time_bucket <= ?", end)
	}
	if !start.IsZero() {
		tx = tx.Where("time_bucket >= ?", start)
	}

	var res []model.BalanceSummary
	return res, tx.Find(&res).Error
}
