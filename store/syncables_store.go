package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

var (
	_ SyncablesStore = (*syncablesStore)(nil)
)

type SyncablesStore interface {
	BaseStore

	FindByHeight(int64) (*model.Syncable, error)
	FindMostRecent() (*model.Syncable, error)
	FindSmallestIndexVersion() (*int64, error)
	FindFirstByDifferentIndexVersion(int64) (*model.Syncable, error)
	FindMostRecentByDifferentIndexVersion(int64) (*model.Syncable, error)
	CreateOrUpdate(*model.Syncable) error
	ResetProcessedAtForRange(int64, int64) error
}

func NewSyncablesStore(db *gorm.DB) *syncablesStore {
	return &syncablesStore{scoped(db, model.Report{})}
}

// syncablesStore handles operations on syncables
type syncablesStore struct {
	baseStore
}

// FindByHeight returns syncable by height
func (s syncablesStore) FindByHeight(height int64) (syncable *model.Syncable, err error) {
	result := &model.Syncable{}

	err = s.db.
		Where("height = ?", height).
		First(result).
		Error

	return result, checkErr(err)
}

// FindMostRecent returns the most recent syncable
func (s syncablesStore) FindMostRecent() (*model.Syncable, error) {
	result := &model.Syncable{}

	err := s.db.
		Order("height desc").
		First(result).Error

	return result, checkErr(err)
}

// FindSmallestIndexVersion returns smallest index version
func (s syncablesStore) FindSmallestIndexVersion() (*int64, error) {
	result := &model.Syncable{}

	err := s.db.
		Where("processed_at IS NOT NULL").
		Order("index_version").
		First(result).Error

	return &result.IndexVersion, checkErr(err)
}

// FindFirstByDifferentIndexVersion returns first syncable with different index version
func (s syncablesStore) FindFirstByDifferentIndexVersion(indexVersion int64) (*model.Syncable, error) {
	result := &model.Syncable{}

	err := s.db.
		Not("index_version = ?", indexVersion).
		Order("height").
		First(result).Error

	return result, checkErr(err)
}

// FindMostRecentByDifferentIndexVersion returns the most recent syncable with different index version
func (s syncablesStore) FindMostRecentByDifferentIndexVersion(indexVersion int64) (*model.Syncable, error) {
	result := &model.Syncable{}

	err := s.db.
		Not("index_version = ?", indexVersion).
		Order("height desc").
		First(result).Error

	return result, checkErr(err)
}

// CreateOrUpdate creates a new syncable or updates an existing one
func (s syncablesStore) CreateOrUpdate(val *model.Syncable) error {
	existing, err := s.FindByHeight(val.Height)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}

	existing.Update(*val)

	return s.Save(existing)
}

// ResetProcessedAtForRange sets processed at to null for given range of heights
func (s syncablesStore) ResetProcessedAtForRange( startHeight int64, endHeight int64) error {
	err := s.db.
		Exec("UPDATE syncables SET processed_at = NULL WHERE height >= ? AND height <= ?", startHeight, endHeight).
		Error

	return checkErr(err)
}
