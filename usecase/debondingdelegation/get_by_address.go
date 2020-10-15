package debondingdelegation

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

type getByAddressUseCase struct {
	db     *store.Store
	client *client.Client
}

func NewGetByAddressUseCase(db *store.Store, c *client.Client) *getByAddressUseCase {
	return &getByAddressUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getByAddressUseCase) Execute(address string, height *int64) (*ListView, error) {
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
		return nil, errors.New("height is not indexed yet")
	}

	res, err := uc.client.DebondingDelegation.GetByAddress(address, *height)
	if err != nil {
		return nil, err
	}

	return ToListViewForAddress(res.GetDebondingDelegations()), nil
}
