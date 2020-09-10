package store

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/jinzhu/gorm"
)

var (
	_ BalanceEventsStore = (*balanceEventsStore)(nil)
)

type BalanceEventsStore interface {
	BaseStore

	CreateOrUpdate(*model.BalanceEvent) error
}

func NewBalanceEventsStore(db *gorm.DB) *balanceEventsStore {
	return &balanceEventsStore{scoped(db, model.BalanceEvent{})}
}

// balanceEventsStore handles operations on syncables
type balanceEventsStore struct {
	baseStore
}

// CreateOrUpdate creates a new system event or updates an existing one
func (s balanceEventsStore) CreateOrUpdate(val *model.BalanceEvent) error {
	existing, err := s.findUnique(val.Height, val.EscrowAddress, val.Address, val.Kind)
	if err != nil {
		if err == ErrNotFound {
			return s.Create(val)
		}
		return err
	}

	existing.Update(*val)

	return s.Save(existing)
}

func (s balanceEventsStore) findUnique(height int64, escrowAddress, address string, kind model.BalanceEventKind) (*model.BalanceEvent, error) {
	q := model.BalanceEvent{
		Height:        height,
		EscrowAddress: escrowAddress,
		Address:       address,
		Kind:          kind,
	}

	var result model.BalanceEvent
	err := s.db.
		Where(&q).
		First(&result).
		Error

	return &result, checkErr(err)
}
