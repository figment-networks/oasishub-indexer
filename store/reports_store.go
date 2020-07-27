package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

var (
	_ ReportsStore = (*reportsStore)(nil)
)

type ReportsStore interface {
	BaseStore

	FindNotCompletedByIndexVersion(int64, ...model.ReportKind) (*model.Report, error)
	FindNotCompletedByKind(...model.ReportKind) (*model.Report, error)
	Last() (*model.Report, error)
	DeleteReindexing() error
}


func NewReportsStore(db *gorm.DB) *reportsStore {
	return &reportsStore{scoped(db, model.Report{})}
}

// reportsStore handles operations on reports
type reportsStore struct {
	baseStore
}

// FindNotCompletedByIndexVersion returns the report by index version and kind
func (s reportsStore) FindNotCompletedByIndexVersion(indexVersion int64, kinds ...model.ReportKind) (*model.Report, error) {
	query := &model.Report{
		IndexVersion: indexVersion,
	}
	result := &model.Report{}

	err := s.db.
		Where(query).
		Where("kind IN(?)", kinds).
		Where("completed_at IS NULL").
		First(result).Error

	return result, checkErr(err)
}

// Last returns the last report
func (s reportsStore) FindNotCompletedByKind(kinds ...model.ReportKind) (*model.Report, error) {
	result := &model.Report{}

	err := s.db.
		Where("kind IN(?)", kinds).
		Where("completed_at IS NULL").
		First(result).Error

	return result, checkErr(err)
}

// Last returns the last report
func (s reportsStore) Last() (*model.Report, error) {
	result := &model.Report{}

	err := s.db.
		Order("id DESC").
		First(result).Error

	return result, checkErr(err)
}

// DeleteReindexing deletes reports with kind reindexing sequential or parallel
func (s *reportsStore) DeleteReindexing() error {
	err := s.db.
		Unscoped().
		Where("kind = ? OR kind = ?", model.ReportKindParallelReindex, model.ReportKindSequentialReindex).
		Delete(&model.Report{}).
		Error

	return checkErr(err)
}