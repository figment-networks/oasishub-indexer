package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

func NewReportsStore(db *gorm.DB) *ReportsStore {
	return &ReportsStore{scoped(db, model.Report{})}
}

// ReportsStore handles operations on reports
type ReportsStore struct {
	baseStore
}

// Last returns the last report
func (s ReportsStore) Last() (*model.Report, error) {
	result := &model.Report{}

	err := s.db.
		Order("id DESC").
		First(result).Error

	return result, checkErr(err)
}
