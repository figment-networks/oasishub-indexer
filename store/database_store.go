package store

import (
	"github.com/jinzhu/gorm"
)

var (
	_ DatabaseStore = (*databaseStore)(nil)
)

type DatabaseStore interface {
	GetTotalSize() (*GetTotalSizeResult, error)
}

func NewDatabaseStore(db *gorm.DB) *databaseStore {
	return &databaseStore{
		db: db,
	}
}

// databaseStore handles operations on blocks
type databaseStore struct {
	db *gorm.DB
}

// GetAvgTimesForIntervalRow Contains row of data for FindSummary query
type GetTotalSizeResult struct {
	Size float64 `json:"size"`
}

// FindSummary Gets average block times for interval
func (s *databaseStore) GetTotalSize() (*GetTotalSizeResult, error) {
	query := "SELECT pg_database_size(current_database()) as size"

	var result GetTotalSizeResult
	err := s.db.Raw(query).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
