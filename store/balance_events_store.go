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
	_ BalanceEventsStore = (*balanceEventsStore)(nil)
)

type BalanceEventsStore interface {
	BaseStore

	CreateOrUpdate(*model.BalanceEvent) error
	Summarize(types.SummaryInterval, []ActivityPeriodRow) ([]model.RawBalanceSummary, error)
}

func NewBalanceEventsStore(db *gorm.DB) *balanceEventsStore {
	return &balanceEventsStore{scoped(db, model.BalanceEvent{})}
}

// balanceEventsStore handles operations on syncables
type balanceEventsStore struct {
	baseStore
}

const (
	summarizeBalanceQuerySelect = `
	s.time_bucket,
	s.min_height,
	balance_events.kind,
	balance_events.address,
	balance_events.escrow_address,
	SUM(balance_events.amount) as total_amount
`
	summarizeBalanceJoinQuery = `INNER JOIN
(
	SELECT
	  MAX(height)     AS max_height,
	  MIN(height)     AS min_height,
	  DATE_TRUNC(?, time) as time_bucket
	FROM syncables
	GROUP BY time_bucket
 ) AS s ON balance_events.height >= s.min_height AND balance_events.height <= s.max_height`
)

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

// Summarize gets the summarized version of balance events
func (s balanceEventsStore) Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]model.RawBalanceSummary, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels("BalanceEventStore_Summarize"))
	defer t.ObserveDuration()

	tx := s.db.
		Table(model.BalanceEvent{}.TableName()).
		Select(summarizeBalanceQuerySelect).
		Joins(summarizeBalanceJoinQuery, interval).
		Group("s.time_bucket, s.min_height, balance_events.address, balance_events.kind, balance_events.escrow_address")

	if len(activityPeriods) == 1 {
		activityPeriod := activityPeriods[0]
		tx = tx.Or("time_bucket < ? OR time_bucket >= ?", activityPeriod.Min, activityPeriod.Max)
	} else {
		for i, activityPeriod := range activityPeriods {
			isLast := i == len(activityPeriods)-1

			if isLast {
				tx = tx.Or("time_bucket >= ?", activityPeriod.Max)
			} else {
				duration, err := time.ParseDuration(fmt.Sprintf("1%s", interval))
				if err != nil {
					return nil, err
				}
				tx = tx.Or("time_bucket >= ? AND time_bucket < ?", activityPeriod.Max.Add(duration), activityPeriods[i+1].Min)
			}
		}
	}

	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []model.RawBalanceSummary
	for rows.Next() {
		var summary model.RawBalanceSummary
		if err := s.db.ScanRows(rows, &summary); err != nil {
			return nil, err
		}

		models = append(models, summary)
	}
	return models, nil
}

func (s balanceEventsStore) findUnique(height int64, escrowAddress, address string, kind model.BalanceEventKind) (*model.BalanceEvent, error) {
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
