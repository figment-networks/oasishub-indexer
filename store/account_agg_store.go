package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

func NewAccountAggStore(db *gorm.DB) *AccountAggStore {
	return &AccountAggStore{scoped(db, model.AccountAgg{})}
}

// AccountAggStore handles operations on accounts
type AccountAggStore struct {
	baseStore
}

// CreateOrUpdate creates a new account or updates an existing one
func (s AccountAggStore) CreateOrUpdate(acc *model.AccountAgg) error {
	existing, err := s.FindByPublicKey(acc.PublicKey)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(acc)
		}
		return err
	}

	return s.Update(existing)
}

// FindBy returns an account for a matching attribute
func (s AccountAggStore) FindBy(key string, value interface{}) (*model.AccountAgg, error) {
	result := &model.AccountAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByID returns an account for the ID
func (s AccountAggStore) FindByID(id int64) (*model.AccountAgg, error) {
	return s.FindBy("id", id)
}

// FindByPublicKey returns an account for the public key
func (s AccountAggStore) FindByPublicKey(key string) (*model.AccountAgg, error) {
	return s.FindBy("public_key", key)
}

// All returns all accounts
func (s AccountAggStore) All() ([]model.AccountAgg, error) {
	var result []model.AccountAgg

	err := s.db.
		Order("id ASC").
		Find(&result).
		Error

	return result, checkErr(err)
}
