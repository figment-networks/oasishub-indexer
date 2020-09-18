package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"
)

var (
	_ BalanceEventsStore = (*balanceEventsStore)(nil)
)

type BalanceEventsStore interface {
	BaseStore

	GetLastEventTime() (types.Time, error)
	CreateOrUpdate(*model.BalanceEvent) error
	DeleteOlderThan(time.Time) (*int64, error)
	Summarize(types.SummaryInterval, []ActivityPeriodRow) ([]model.BalanceSummary, error)
}

func NewBalanceEventsStore(db *gorm.DB) *balanceEventsStore {
	return &balanceEventsStore{scoped(db, model.BalanceEvent{})}
}

// balanceEventsStore handles operations on syncables
type balanceEventsStore struct {
	baseStore
}

// CreateOrUpdate creates a new balance event or updates an existing one
func (s *balanceEventsStore) CreateOrUpdate(val *model.BalanceEvent) error {
	existing, err := s.findUnique(val.Height, val.EscrowAddress, val.Address, val.Kind)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}

	existing.Update(*val)
	return s.Save(existing)
}

// DeleteOlderThan deletes balance events older than given threshold
func (s *balanceEventsStore) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	query := s.db.Table("syncables").Select("height").Where("time < ?", purgeThreshold).QueryExpr()

	tx := s.db.
		Unscoped().
		Where("height IN (?)", query).
		Delete(&model.BalanceEvent{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}

// GetLastEventTime returns the time corresponding to the most recent balance event
func (s *balanceEventsStore) GetLastEventTime() (types.Time, error) {
	var result struct {
		Time types.Time `json:"time"`
	}

	err := s.db.
		Table(model.BalanceEvent{}.TableName()).
		Select("s.time").
		Joins("INNER JOIN syncables AS s ON balance_events.height = s.height").
		Order("s.time DESC").
		First(&result).
		Error

	return result.Time, checkErr(err)
}

// Summarize gets the summarized version of balance events
func (s *balanceEventsStore) Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]model.BalanceSummary, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels("BalanceEventStore_Summarize"))
	defer t.ObserveDuration()

	tx := s.db.
		Table(model.BalanceEvent{}.TableName()).
		Select(summarizeBalanceQuerySelect).
		Joins(summarizeBalanceJoinQuery, interval).
		Group("s.time_bucket, balance_events.address, balance_events.escrow_address, s.start_height")

	if len(activityPeriods) == 1 {
		activityPeriod := activityPeriods[0]
		tx = tx.Or("time_bucket < ? OR time_bucket >= ?", activityPeriod.Min, activityPeriod.Max)
	} else {
		for i, activityPeriod := range activityPeriods {
			isLast := i == len(activityPeriods)-1

			if isLast {
				tx = tx.Or("time_bucket >= ?", activityPeriod.Max)
			} else {
				duration, err := interval.ToDuration()
				if err != nil {
					return nil, err
				}
				tx = tx.Or("time_bucket >= ? AND time_bucket < ?", activityPeriod.Max.Add(duration), activityPeriods[i+1].Min)
			}
		}
	}

	var models []model.BalanceSummary
	return models, tx.Find(&models).Error
}

func (s *balanceEventsStore) findUnique(height int64, escrowAddress, address string, kind model.BalanceEventKind) (*model.BalanceEvent, error) {
	q := model.BalanceEvent{
		Height:        height,
		EscrowAddress: escrowAddress,
		Address:       address,
		Kind:          kind,
	}

	var result model.BalanceEvent
	err := s.db.
		Where(&q).
		First(&result).
		Error

	return &result, checkErr(err)
}
