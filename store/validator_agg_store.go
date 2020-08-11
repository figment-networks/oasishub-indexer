package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

var (
	_ ValidatorAggStore = (*validatorAggStore)(nil)
)

type ValidatorAggStore interface {
	BaseStore

	FindBy(string, interface{}) (*model.ValidatorAgg, error)
	FindByAddress(string) (*model.ValidatorAgg, error)
	FindByEntityUID(string) (*model.ValidatorAgg, error)
	GetAllForHeightGreaterThan(int64) ([]model.ValidatorAgg, error)
	CreateOrUpdate(val *model.ValidatorAgg) error
}


func NewValidatorAggStore(db *gorm.DB) *validatorAggStore {
	return &validatorAggStore{scoped(db, model.ValidatorAgg{})}
}

// validatorAggStore handles operations on validators
type validatorAggStore struct {
	baseStore
}

// FindBy returns an validator for a matching attribute
func (s validatorAggStore) FindBy(key string, value interface{}) (*model.ValidatorAgg, error) {
	result := &model.ValidatorAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByAddress return validator by entity UID
func (s *validatorAggStore) FindByAddress(address string) (*model.ValidatorAgg, error) {
	return s.FindBy("address", address)
}

// FindByEntityUID return validator by entity UID
func (s *validatorAggStore) FindByEntityUID(key string) (*model.ValidatorAgg, error) {
	return s.FindBy("entity_uid", key)
}

// GetAllForHeightGreaterThan returns validators who have been validating since given height
func (s *validatorAggStore) GetAllForHeightGreaterThan(height int64) ([]model.ValidatorAgg, error) {
	var result []model.ValidatorAgg

	err := s.baseStore.db.
		Where("recent_as_validator_height >= ?", height).
		Find(&result).
		Error

	return result, checkErr(err)
}

// CreateOrUpdate creates a new validator or updates an existing one
func (s validatorAggStore) CreateOrUpdate(val *model.ValidatorAgg) error {
	_, err := s.FindByEntityUID(val.EntityUID)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}
	return s.Update(val)
}