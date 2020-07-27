package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

var (
	_ AccountAggStore = (*accountAggStore)(nil)
)

type AccountAggStore interface {
	BaseStore

	FindBy(string, interface{}) (*model.AccountAgg, error)
	FindByPublicKey(string) (*model.AccountAgg, error)
}

func NewAccountAggStore(db *gorm.DB) *accountAggStore {
	return &accountAggStore{scoped(db, model.AccountAgg{})}
}

// accountAggStore handles operations on accounts
type accountAggStore struct {
	baseStore
}

// FindBy returns an account for a matching attribute
func (s accountAggStore) FindBy(key string, value interface{}) (*model.AccountAgg, error) {
	result := &model.AccountAgg{}
	err := findBy(s.db, result, key, value)
	return result, checkErr(err)
}

// FindByPublicKey returns an account for the public key
func (s accountAggStore) FindByPublicKey(key string) (*model.AccountAgg, error) {
	return s.FindBy("public_key", key)
}
