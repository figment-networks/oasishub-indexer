package syncable

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getMostRecentHeightUseCase struct {
	db *store.Store
}

func NewGetMostRecentHeightUseCase(db *store.Store) *getMostRecentHeightUseCase {
	return &getMostRecentHeightUseCase{
		db: db,
	}
}

func (uc *getMostRecentHeightUseCase) Execute() (*int64, error) {
	h, err := uc.db.Syncables.FindMostRecent()
	if err != nil {
		return nil, err
	}
	return &h.Height, nil
}

