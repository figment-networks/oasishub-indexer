package validator

import (
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

type getForMinHeightUseCase struct {
	db *store.Store
}

func NewGetForMinHeightUseCase(db *store.Store) *getForMinHeightUseCase {
	return &getForMinHeightUseCase{
		db: db,
	}
}

func (uc *getForMinHeightUseCase) Execute(height *int64) (*AggListView, error) {
	// Get last indexed height
	mostRecentSynced, err := uc.db.Syncables.FindMostRecent()
	if err != nil {
		return nil, err
	}
	lastH := mostRecentSynced.Height

	// Show last synced height, if not provided
	if height == nil {
		height = &lastH
	}

	if *height > lastH {
		return nil, errors.New("height is not indexed")
	}

	ms, err := uc.db.ValidatorAgg.GetAllForHeightGreaterThan(*height)
	if err != nil {
		return nil, err
	}

	return ToAggListView(ms), nil
}


