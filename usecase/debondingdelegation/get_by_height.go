package debondingdelegation

import (
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

type getByHeightUseCase struct {
	db *store.Store
}

func NewGetByHeightUseCase(db *store.Store) *getByHeightUseCase {
	return &getByHeightUseCase{
		db: db,
	}
}

func (uc *getByHeightUseCase) Execute(height *int64) (*ListView, error) {
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

	ms, err := uc.db.DebondingDelegationSeq.FindByHeight(*height)
	if err != nil {
		return nil, err
	}

	return ToListView(ms), nil
}
