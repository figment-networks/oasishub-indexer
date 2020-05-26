package store

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"

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

// GetTotalValidatedByEntityUID gets total validated blocks for validator
func (s *ValidatorSeqStore) GetTotalValidatedByEntityUID(key string) (*int64, error) {
	v := true
	q := model.ValidatorSeq{
		EntityUID:          key,
		PrecommitValidated: &v,
	}
	var result int64

	err := s.db.
		Table(model.ValidatorSeq{}.TableName()).
		Where(&q).
		Count(&result).
		Error

	return &result, checkErr(err)
}

// GetTotalMissedByEntityUID gets total missed blocks for validator
func (s *ValidatorSeqStore) GetTotalMissedByEntityUID(key string) (*int64, error) {
	v := false
	q := model.ValidatorSeq{
		EntityUID:          key,
		PrecommitValidated: &v,
	}
	var result int64

	err := s.db.
		Table(model.ValidatorSeq{}.TableName()).
		Where(&q).
		Count(&result).
		Error

	return &result, checkErr(err)
}

// GetTotalProposedByEntityUID gets total proposed blocks for validator
func (s *ValidatorSeqStore) GetTotalProposedByEntityUID(key string) (*int64, error) {
	q := model.ValidatorSeq{
		EntityUID: key,
		Proposed:  true,
	}
	var result int64

	err := s.db.
		Table(model.ValidatorSeq{}.TableName()).
		Where(&q).
		Count(&result).
		Error

	return &result, checkErr(err)
}

// AvgForTimeIntervalRow result of time interval query
type AvgForTimeIntervalRow struct {
	TimeInterval string  `json:"time_interval"`
	Avg          float64 `json:"avg"`
}

// AvgQuantityForTimeIntervalRow result of time interval query with quantity
type AvgQuantityForTimeIntervalRow struct {
	TimeInterval string         `json:"time_interval"`
	Avg          types.Quantity `json:"avg"`
}

// GetTotalSharesForInterval gets total shares of all validators for interval
func (s *ValidatorSeqStore) GetTotalSharesForInterval(interval string, period string) ([]AvgQuantityForTimeIntervalRow, error) {
	rows, err := s.db.Raw(totalSharesForIntervalQuery, interval, period).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AvgQuantityForTimeIntervalRow
	for rows.Next() {
		var row AvgQuantityForTimeIntervalRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// GetTotalVotingPowerForInterval gets total voting power of all validators for interval
func (s *ValidatorSeqStore) GetTotalVotingPowerForInterval(interval string, period string) ([]AvgForTimeIntervalRow, error) {
	rows, err := s.db.Raw(totalVotingPowerForIntervalQuery, interval, period).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AvgForTimeIntervalRow
	for rows.Next() {
		var row AvgForTimeIntervalRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// GetValidatorSharesForInterval gets shares for validator for interval
func (s *ValidatorSeqStore) GetValidatorSharesForInterval(key string, interval string, period string) ([]AvgQuantityForTimeIntervalRow, error) {
	rows, err := s.db.Raw(validatorSharesForIntervalQuery, key, interval, period).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AvgQuantityForTimeIntervalRow
	for rows.Next() {
		var row AvgQuantityForTimeIntervalRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// GetValidatorVotingPowerForInterval gets voting powers for validator for interval
func (s *ValidatorSeqStore) GetValidatorVotingPowerForInterval(key string, interval string, period string) ([]AvgForTimeIntervalRow, error) {
	rows, err := s.db.Debug().Raw(validatorVotingPowerForIntervalQuery, key, interval, period).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AvgForTimeIntervalRow
	for rows.Next() {
		var row AvgForTimeIntervalRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// GetValidatorUptimeForInterval gets uptime for validator for interval
func (s *ValidatorSeqStore) GetValidatorUptimeForInterval(key string, interval string, period string) ([]AvgForTimeIntervalRow, error) {
	rows, err := s.db.Raw(validatorUptimeForIntervalQuery, key, interval, period).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AvgForTimeIntervalRow
	for rows.Next() {
		var row AvgForTimeIntervalRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}