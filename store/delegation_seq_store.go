package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

func NewDelegationSeqStore(db *gorm.DB) *DelegationSeqStore {
	return &DelegationSeqStore{scoped(db, model.DelegationSeq{})}
}

// DelegationSeqStore handles operations on delegations
type DelegationSeqStore struct {
	baseStore
}

// CreateIfNotExists creates the delegation if it does not exist
func (s DelegationSeqStore) CreateIfNotExists(delegation *model.DelegationSeq) error {
	_, err := s.FindByHeight(delegation.Height)
	if isNotFound(err) {
		return s.Create(delegation)
	}
	return nil
}

// FindByHeight finds delegation by height
func (s DelegationSeqStore) FindByHeight(h int64) ([]model.DelegationSeq, error) {
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
func (s *DelegationSeqStore) FindLastByValidatorUID(key string) ([]model.DelegationSeq, error) {
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
func (s *DelegationSeqStore) FindCurrentByDelegatorUID(key string) ([]model.DelegationSeq, error) {
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
