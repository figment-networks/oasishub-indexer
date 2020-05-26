package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/oasishub-indexer/model"
)

func NewBlockSeqStore(db *gorm.DB) *BlockSeqStore {
	return &BlockSeqStore{scoped(db, model.BlockSeq{})}
}

// BlockSeqStore handles operations on blocks
type BlockSeqStore struct {
	baseStore
}

// CreateIfNotExists creates the block if it does not exist
func (s BlockSeqStore) CreateIfNotExists(block *model.BlockSeq) error {
	_, err := s.FindByHeight(block.Height)
	if isNotFound(err) {
		return s.Create(block)
	}
	return nil
}

// FindBy returns a block for a matching attribute
func (s BlockSeqStore) FindBy(key string, value interface{}) (*model.BlockSeq, error) {
	result := &model.BlockSeq{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns a block with matching ID
func (s BlockSeqStore) FindByID(id int64) (*model.BlockSeq, error) {
	return s.FindBy("id", id)
}

// FindByHeight returns a block with the matching height
func (s BlockSeqStore) FindByHeight(height int64) (*model.BlockSeq, error) {
	return s.FindBy("height", height)
}

// GetAvgRecentTimesResult Contains results for GetAvgRecentTimes query
type GetAvgRecentTimesResult struct {
	StartHeight int64   `json:"start_height"`
	EndHeight   int64   `json:"end_height"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Count       int64   `json:"count"`
	Diff        float64 `json:"diff"`
	Avg         float64 `json:"avg"`
}

// GetAvgRecentTimes Gets average block times for recent blocks by limit
func (s *BlockSeqStore) GetAvgRecentTimes(limit int64) GetAvgRecentTimesResult {
	var res GetAvgRecentTimesResult
	s.db.Raw(blockTimesForRecentBlocksQuery, limit).Scan(&res)

	return res
}

// GetAvgTimesForIntervalRow Contains row of data for GetAvgTimesForInterval query
type GetAvgTimesForIntervalRow struct {
	TimeInterval string  `json:"time_interval"`
	Count        int64   `json:"count"`
	Avg          float64 `json:"avg"`
}

// GetAvgTimesForInterval Gets average block times for interval
func (s *BlockSeqStore) GetAvgTimesForInterval(interval string, period string) ([]GetAvgTimesForIntervalRow, error) {
	rows, err := s.db.Raw(blockTimesForIntervalQuery, interval, period).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []GetAvgTimesForIntervalRow
	for rows.Next() {
		var row GetAvgTimesForIntervalRow
		if err := s.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, nil
}