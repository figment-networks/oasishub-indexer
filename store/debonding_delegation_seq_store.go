package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/oasishub-indexer/model"
)

var (
	_ DebondingDelegationSeqStore = (*debondingDelegationSeqStore)(nil)
)

type DebondingDelegationSeqStore interface {
	BaseStore

	FindByHeight(int64) ([]model.DebondingDelegationSeq, error)
	FindRecentByValidatorUID(string, int64) ([]model.DebondingDelegationSeq, error)
	FindRecentByDelegatorUID(string, int64) ([]model.DebondingDelegationSeq, error)
}

func NewDebondingDelegationSeqStore(db *gorm.DB) *debondingDelegationSeqStore {
	return &debondingDelegationSeqStore{scoped(db, model.DebondingDelegationSeq{})}
}

// debondingDelegationSeqStore handles operations on debondingDelegations
type debondingDelegationSeqStore struct {
	baseStore
}

// FindByHeight finds debonding delegations by height
func (s debondingDelegationSeqStore) FindByHeight(h int64) ([]model.DebondingDelegationSeq, error) {
	q := model.DebondingDelegationSeq{
		Sequence: &model.Sequence{
			Height: h,
		},
	}
	var result []model.DebondingDelegationSeq

	err := s.db.Where(&q).Find(&result).Error
	return result, checkErr(err)
}

// FindRecentByValidatorUID gets recent debonding delegations for validator
func (s *debondingDelegationSeqStore) FindRecentByValidatorUID(key string, limit int64) ([]model.DebondingDelegationSeq, error) {
	q := model.DebondingDelegationSeq{
		ValidatorUID:  key,
	}
	var result []model.DebondingDelegationSeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

// FindRecentByDelegatorUID gets recent debonding delegations for delegator
func (s *debondingDelegationSeqStore) FindRecentByDelegatorUID(key string, limit int64) ([]model.DebondingDelegationSeq, error) {
	q := model.DebondingDelegationSeq{
		DelegatorUID:  key,
	}
	var result []model.DebondingDelegationSeq

	err := s.db.
		Where(&q).
		Order("height DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}