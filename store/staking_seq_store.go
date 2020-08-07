package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/oasishub-indexer/model"
)

var (
	_ StakingSeqStore = (*stakingSeqStore)(nil)
)

type StakingSeqStore interface {
	BaseStore

	FindBy(key string, value interface{}) (*model.StakingSeq, error)
	FindByHeight(height int64) (*model.StakingSeq, error)
	Recent() (*model.StakingSeq, error)
}


func NewStakingSeqStore(db *gorm.DB) *stakingSeqStore {
	return &stakingSeqStore{scoped(db, model.StakingSeq{})}
}

// stakingSeqStore handles operations on staking
type stakingSeqStore struct {
	baseStore
}

// FindBy returns a staking for a matching attribute
func (s stakingSeqStore) FindBy(key string, value interface{}) (*model.StakingSeq, error) {
	result := &model.StakingSeq{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByHeight returns a staking with the matching height
func (s stakingSeqStore) FindByHeight(height int64) (*model.StakingSeq, error) {
	return s.FindBy("height", height)
}

// Recent returns the most recent staking
func (s stakingSeqStore) Recent() (*model.StakingSeq, error) {
	staking := &model.StakingSeq{}

	err := s.db.
		Order("height DESC").
		First(staking).
		Error

	return staking, checkErr(err)
}

