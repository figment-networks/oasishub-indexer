package store

import (
	"github.com/figment-networks/oasishub-indexer/config"
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

type ValidatorSummaryRow struct {
	TimeInterval    string         `json:"time_interval"`
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

type SingleValidatorSummaryRow struct {
	EntityUID string `json:"entity_uid"`

	*ValidatorSummaryRow
}

// GetSummary gets total shares of all validators for interval
func (s *ValidatorSeqStore) GetSummary(interval string, period string) ([]ValidatorSummaryRow, error) {
	defer logQueryDuration(time.Now(), "ValidatorSeqStore_GetSummary")

	rows, err := s.db.Raw(AllValidatorsSummaryForIntervalQuery(interval), period).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ValidatorSummaryRow
	for rows.Next() {
		var row ValidatorSummaryRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}


// GetSummaryByEntityUID gets shares for validator for interval
func (s *ValidatorSeqStore) GetSummaryByEntityUID(key string, interval string, period string) ([]SingleValidatorSummaryRow, error) {
	defer logQueryDuration(time.Now(), "ValidatorSeqStore_GetSummaryByEntityUID")

	rows, err := s.db.Raw(ValidatorSummaryForIntervalQuery(interval), period, key).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []SingleValidatorSummaryRow
	for rows.Next() {
		var row SingleValidatorSummaryRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

func (s *ValidatorSeqStore) PurgeOldRecords(cfg *config.Config) error {
	// Purge sequences
	if err := s.db.Exec(deleteOldValidatorSeqQuery, cfg.PurgeValidatorInterval).Error; err != nil {
		return err
	}

	// Purge hourly summary
	if err := s.db.Exec(DeleteOldValidatorHourlySummaryQuery("hourly"), cfg.PurgeValidatorHourlySummaryInterval).Error; err != nil {
		return err
	}

	// Purge daily summary
	if err := s.db.Exec(DeleteOldValidatorHourlySummaryQuery("daily"), cfg.PurgeValidatorDailySummaryInterval).Error; err != nil {
		return err
	}

	return nil
}