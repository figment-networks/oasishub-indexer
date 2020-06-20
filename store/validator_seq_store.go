package store

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/figment-networks/oasishub-indexer/model"
)

func NewValidatorSeqStore(db *gorm.DB) *ValidatorSeqStore {
	return &ValidatorSeqStore{scoped(db, model.ValidatorSeq{})}
}

// ValidatorSeqStore handles operations on validators
type ValidatorSeqStore struct {
	baseStore
}

// CreateIfNotExists creates the validator if it does not exist
func (s ValidatorSeqStore) CreateIfNotExists(validator *model.ValidatorSeq) error {
	_, err := s.FindByHeight(validator.Height)
	if isNotFound(err) {
		return s.Create(validator)
	}
	return nil
}

//TODO: FIX
// CreateOrUpdate creates a new validator sequence or updates an existing one
func (s ValidatorSeqStore) CreateOrUpdate(val *model.ValidatorSeq) error {
	existing, err := s.FindByHeight(val.Height)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}
	return s.Update(existing)
}

// FindByHeight finds validator by height
func (s ValidatorSeqStore) FindByHeightAndEntityUID(h int64, key string) (*model.ValidatorSeq, error) {
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
func (s ValidatorSeqStore) FindByHeight(h int64) ([]model.ValidatorSeq, error) {
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
func (s ValidatorSeqStore) FindLastByAddress(address string, limit int64) ([]model.ValidatorSeq, error) {
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
func (s *ValidatorSeqStore) FindMostRecent() (*model.ValidatorSeq, error) {
	validatorSeq := &model.ValidatorSeq{}
	if err := findMostRecent(s.db, "time", validatorSeq); err != nil {
		return nil, err
	}
	return validatorSeq, nil
}

// DeleteOlderThan deletes validator sequence older than given threshold
func (s *ValidatorSeqStore) DeleteOlderThan(purgeThreshold time.Time) (*int64, error) {
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
	Address         string         `json:"address"`
	TimeBucket      types.Time     `json:"time_bucket"`
	VotingPowerAvg  float64        `json:"voting_power_avg"`
	VotingPowerMax  float64        `json:"voting_power_max"`
	VotingPowerMin  float64        `json:"voting_power_min"`
	TotalSharesAvg  types.Quantity `json:"total_shares_avg"`
	TotalSharesMax  types.Quantity `json:"total_shares_max"`
	TotalSharesMin  types.Quantity `json:"total_shares_min"`
	ValidatedSum    int64          `json:"validated_sum"`
	NotValidatedSum int64          `json:"not_validated_sum"`
	ProposedSum     int64          `json:"proposed_sum"`
	UptimeAvg       float64        `json:"uptime_avg"`
}

// Summarize gets the summarized version of validator sequences
func (s *ValidatorSeqStore) Summarize(interval types.SummaryInterval, activityPeriods []ActivityPeriodRow) ([]ValidatorSeqSummary, error) {
	defer logQueryDuration(time.Now(), "ValidatorSeqStore_Summarize")

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

	var models []ValidatorSeqSummary
	for rows.Next() {
		var summary ValidatorSeqSummary
		if err := s.db.ScanRows(rows, &summary); err != nil {
			return nil, err
		}

		models = append(models, summary)
	}
	return models, nil
}
