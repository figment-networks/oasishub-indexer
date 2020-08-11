package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

var (
	_ DelegationSeqStore = (*delegationSeqStore)(nil)
)

type DelegationSeqStore interface {
	BaseStore

	FindByHeight(int64) ([]model.DelegationSeq, error)
	FindLastByValidatorUID(string) ([]model.DelegationSeq, error)
	FindCurrentByDelegatorUID(string) ([]model.DelegationSeq, error)
}

func NewDelegationSeqStore(db *gorm.DB) *delegationSeqStore {
	return &delegationSeqStore{scoped(db, model.DelegationSeq{})}
}

// delegationSeqStore handles operations on delegations
type delegationSeqStore struct {
	baseStore
}

// FindByHeight finds delegation by height
func (s delegationSeqStore) FindByHeight(h int64) ([]model.DelegationSeq, error) {
	q := model.DelegationSeq{
		Sequence: &model.Sequence{
			Height: h,
		},
	}
	var result []model.DelegationSeq

	err := s.db.Where(&q).Find(&result).Error
	return result, checkErr(err)
}

// GetLastByValidatorUID finds last delegations for validator
func (s *delegationSeqStore) FindLastByValidatorUID(key string) ([]model.DelegationSeq, error) {
	q := model.DelegationSeq{
		ValidatorUID:  key,
	}
	var result []model.DelegationSeq

	sub := s.db.Table(model.DelegationSeq{}.TableName()).Select("height").Order("height DESC").Limit(1).QueryExpr()
	err := s.db.
		Where(&q).
		Where("height = (?)", sub).
		Find(&result).
		Error

	return result, checkErr(err)
}

// GetCurrentByDelegatorUID gets current delegations for delegator
func (s *delegationSeqStore) FindCurrentByDelegatorUID(key string) ([]model.DelegationSeq, error) {
	q := model.DelegationSeq{
		DelegatorUID:  key,
	}
	var result []model.DelegationSeq

	sub := s.db.Table(model.DelegationSeq{}.TableName()).Select("height").Order("height DESC").Limit(1).QueryExpr()
	err := s.db.
		Where(&q).
		Where("height = (?)", sub).
		Find(&result).
		Error

	return result, checkErr(err)
}
