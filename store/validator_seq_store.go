package store

import (
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

// FindLastByEntityUID finds last validator sequences for given entity uid
func (s ValidatorSeqStore) FindLastByEntityUID(key string, limit int64) ([]model.ValidatorSeq, error) {
	q := model.ValidatorSeq{
		EntityUID: key,
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
	statement := s.db.
		Unscoped().
		Where("time < ?", purgeThreshold).
		Delete(&model.ValidatorSeq{})

	if statement.Error != nil {
		return nil, checkErr(statement.Error)
	}

	return &statement.RowsAffected, nil
}

type ValidatorSeqSummary struct {
	EntityUID       string         `json:"entity_uid"`
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
func (s *ValidatorSeqStore) Summarize(interval types.SummaryInterval, last *model.ValidatorSummary) ([]ValidatorSeqSummary, error) {
	defer logQueryDuration(time.Now(), "ValidatorSeqStore_Summarize")

	tx := s.db.
		Table(model.ValidatorSeq{}.TableName()).
		Select(summarizeValidatorsQuerySelect, interval).
		Order("time_bucket").
		Group("entity_uid, time_bucket")

	if last != nil {
		tx = tx.Where("time >= ?", last.TimeBucket)
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
