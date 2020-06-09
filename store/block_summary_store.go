package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/jinzhu/gorm"
	"time"
)

func NewBlockSummaryStore(db *gorm.DB) *BlockSummaryStore {
	return &BlockSummaryStore{scoped(db, model.BlockSummary{})}
}

// BlockSummaryStore handles operations on block summary
type BlockSummaryStore struct {
	baseStore
}

// Find find block summary by query
func (s BlockSummaryStore) Find(query *model.BlockSummary) (*model.BlockSummary, error) {
	var result model.BlockSummary

	err := s.db.
		Where(query).
		First(&result).
		Error

	return &result, checkErr(err)
}

// FindMostRecent finds most recent block summary
func (s *BlockSummaryStore) FindMostRecent() (*model.BlockSummary, error) {
	blockSummary := &model.BlockSummary{}
	err := findMostRecent(s.db, "time_bucket", blockSummary)
	return blockSummary, checkErr(err)
}

// FindMostRecentByInterval finds most recent block summary for given time interval
func (s *BlockSummaryStore) FindMostRecentByInterval(interval types.SummaryInterval) (*model.BlockSummary, error) {
	query := &model.BlockSummary{
		Summary: &model.Summary{TimeInterval: interval},
	}
	result := model.BlockSummary{}

	err := s.db.
		Where(query).
		Order("time_bucket DESC").
		Take(&result).
		Error

	return &result, checkErr(err)
}

// FindSummary Gets summary of block sequences
func (s *BlockSummaryStore) FindSummary(interval types.SummaryInterval, period string) ([]model.BlockSummary, error) {
	defer logQueryDuration(time.Now(), "BlockSummaryStore_FindSummary")

	rows, err := s.db.
		Raw(allBlocksSummaryForIntervalQuery, interval, period, interval).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.BlockSummary
	for rows.Next() {
		var row model.BlockSummary
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}

// DeleteOlderThan deletes block summary records older than given threshold
func (s *BlockSummaryStore) DeleteOlderThan(interval types.SummaryInterval, purgeThreshold time.Time) (*int64, error) {
	res := s.db.
		Unscoped().
		Where("time_interval = ? AND time_bucket < ?", interval, purgeThreshold).
		Delete(&model.BlockSummary{})

	if res.Error != nil {
		return nil, checkErr(res.Error)
	}

	return &res.RowsAffected, nil
}