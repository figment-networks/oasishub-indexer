package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/oasishub-indexer/model"
)

func NewStakingSeqStore(db *gorm.DB) *StakingSeqStore {
	return &StakingSeqStore{scoped(db, model.StakingSeq{})}
}

// StakingSeqStore handles operations on staking
type StakingSeqStore struct {
	baseStore
}

// CreateIfNotExists creates the staking if it does not exist
func (s StakingSeqStore) CreateIfNotExists(staking *model.StakingSeq) error {
	_, err := s.FindByHeight(staking.Height)
	if isNotFound(err) {
		return s.Create(staking)
	}
	return nil
}

// FindBy returns a staking for a matching attribute
func (s StakingSeqStore) FindBy(key string, value interface{}) (*model.StakingSeq, error) {
	result := &model.StakingSeq{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns a staking with matching ID
func (s StakingSeqStore) FindByID(id int64) (*model.StakingSeq, error) {
	return s.FindBy("id", id)
}

// FindByHeight returns a staking with the matching height
func (s StakingSeqStore) FindByHeight(height int64) (*model.StakingSeq, error) {
	return s.FindBy("height", height)
}

// Recent returns the most recent staking
func (s StakingSeqStore) Recent() (*model.StakingSeq, error) {
	staking := &model.StakingSeq{}

	err := s.db.
		Order("height DESC").
		First(staking).
		Error

	return staking, checkErr(err)
}
