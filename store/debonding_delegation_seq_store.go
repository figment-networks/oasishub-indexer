package store

import (
	"github.com/jinzhu/gorm"

	"github.com/figment-networks/oasishub-indexer/model"
)

func NewDebondingDelegationSeqStore(db *gorm.DB) *DebondingDelegationSeqStore {
	return &DebondingDelegationSeqStore{scoped(db, model.DebondingDelegationSeq{})}
}

// DebondingDelegationSeqStore handles operations on debondingDelegations
type DebondingDelegationSeqStore struct {
	baseStore
}

// CreateIfNotExists creates the debondingDelegation if it does not exist
func (s DebondingDelegationSeqStore) CreateIfNotExists(debondingDelegation *model.DebondingDelegationSeq) error {
	_, err := s.FindByHeight(debondingDelegation.Height)
	if isNotFound(err) {
		return s.Create(debondingDelegation)
	}
	return nil
}

// FindByHeight finds debonding delegations by height
func (s DebondingDelegationSeqStore) FindByHeight(h int64) ([]model.DebondingDelegationSeq, error) {
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
func (s *DebondingDelegationSeqStore) FindRecentByValidatorUID(key string, limit int64) ([]model.DebondingDelegationSeq, error) {
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
func (s *DebondingDelegationSeqStore) FindRecentByDelegatorUID(key string, limit int64) ([]model.DebondingDelegationSeq, error) {
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