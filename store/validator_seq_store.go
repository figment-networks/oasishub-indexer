package store

import (
	"time"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"
)

var (
	_ ValidatorSeqStore = (*validatorSeqStore)(nil)
)

type ValidatorSeqStore interface {
	BaseStore

	FindByHeightAndEntityUID(int64, string) (*model.ValidatorSeq, error)
	FindByHeight(int64) ([]model.ValidatorSeq, error)
	FindLastByAddress(string, int64) ([]model.ValidatorSeq, error)
	FindMostRecent() (*model.ValidatorSeq, error)
	DeleteOlderThan(time.Time) (*int64, error)
	Summarize(types.SummaryInterval, []ActivityPeriodRow) ([]ValidatorSeqSummary, error)
}

func NewValidatorSeqStore(db *gorm.DB) *validatorSeqStore {
	return &validatorSeqStore{scoped(db, model.ValidatorSeq{})}
}

// validatorSeqStore handles operations on validators
type validatorSeqStore struct {
	baseStore
}

// FindByHeightAndEntityUID finds validator by height amd entity UID
func (s validatorSeqStore) FindByHeightAndEntityUID(h int64, key string) (*model.ValidatorSeq, error) {
	q := model.ValidatorSeq{
		Sequence: &model.Sequence{
			Height: h,
		},
		EntityUID: key,
	}
	var result model.ValidatorSeq

	err := s.db.
		Where(&q).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindByHeight finds validator by height
func (s validatorSeqStore) FindByHeight(h int64) ([]model.ValidatorSeq, error) {
	q := model.ValidatorSeq{
		Sequence: &model.Sequence{
			Height: h,
		},
	}
	var result []model.ValidatorSeq

	err := s.db.
		Where(&q).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindLastByAddress finds last validator sequences for given entity uid
func (s validatorSeqStore) FindLastByAddress(address string, limit int64) ([]model.ValidatorSeq, error) {
	q := model.ValidatorSeq{
		Address: address,
	}
	var result []model.ValidatorSeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindMostRecent finds most recent validator sequence
func (s *validatorSeqStore) FindMostRecent() (*model.ValidatorSeq, error) {
	validatorSeq := &model.ValidatorSeq{}
	if err := findMostRecent(s.db, "time", validatorSeq); err != nil {
		return nil, err
	}
	return validatorSeq, nil
}

// DeleteOlderThan deletes validator sequence older than given threshold
func (s *validatorSeqStore) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
	tx := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.ValidatorSeq{})

	if tx.Error != nil {
		return nil, checkErr(tx.Error)
	}

	return &tx.RowsAffected, nil
}

type ValidatorSeqSummary struct {
	Address                string         `json:"address"`
	TimeBucket             types.Time     `json:"time_bucket"`
	VotingPowerAvg         float64        `json:"voting_power_avg"`
	VotingPowerMax         float64        `json:"voting_power_max"`
	VotingPowerMin         float64        `json:"voting_power_min"`
	TotalSharesAvg         types.Quantity `json:"total_shares_avg"`
	TotalSharesMax         types.Quantity `json:"total_shares_max"`
	TotalSharesMin         types.Quantity `json:"total_shares_min"`
	ActiveEscrowBalanceAvg types.Quantity `json:"active_escrow_balance_avg"`
	ActiveEscrowBalanceMax types.Quantity `json:"active_escrow_balance_max"`
	ActiveEscrowBalanceMin types.Quantity `json:"active_escrow_balance_min"`
	CommissionAvg          types.Quantity `json:"commission_avg"`
	CommissionMax          types.Quantity `json:"commission_max"`
	CommissionMin          types.Quantity `json:"commission_min"`
	ValidatedSum           int64          `json:"validated_sum"`
	NotValidatedSum        int64          `json:"not_validated_sum"`
	ProposedSum            int64          `json:"proposed_sum"`
	UptimeAvg              float64        `json:"uptime_avg"`
}

// Summarize gets the summarized version of validator sequences
func (s *validatorSeqStore) Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]ValidatorSeqSummary, error) {
	t := metrics.NewTimer(databaseQueryDuration.WithLabels("ValidatorSeqStore_Summarize"))
	defer t.ObserveDuration()

	tx := s.db.
		Table(model.ValidatorSeq{}.TableName()).
		Select(summarizeValidatorsQuerySelect, interval).
		Order("time_bucket").
		Group("address, time_bucket")

	if len(activityPeriods) == 1 {
		activityPeriod := activityPeriods[0]
		tx = tx.Or("time < ? OR time >= ?", activityPeriod.Min, activityPeriod.Max)
	} else {
		for i, activityPeriod := range activityPeriods {
			isLast := i == len(activityPeriods)-1

			if isLast {
				tx = tx.Or("time >= ?", activityPeriod.Max)
			} else {
				duration, err := interval.ToDuration()
				if err != nil {
					return nil, err
				}
				tx = tx.Or("time >= ? AND time < ?", activityPeriod.Max.Add(duration), activityPeriods[i+1].Min)
			}
		}
	}

	var models []ValidatorSeqSummary
	return models, tx.Find(&models).Error
}
