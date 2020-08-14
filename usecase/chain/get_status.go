package chain

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
)

type getStatusUseCase struct {
	db     *store.Store
	client *client.Client
}

func NewGetStatusUseCase(db *store.Store, c *client.Client) *getStatusUseCase {
	return &getStatusUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getStatusUseCase) Execute(ctx  context.Context) (*DetailsView, error) {
	mostRecentSyncable, err := uc.db.Syncables.FindMostRecent()
	if err != nil {
		return nil, err
	}

	getHeadRes, err := uc.client.Chain.GetHead()
	if err != nil {
		return nil, err
	}

	getStatusRes, err := uc.client.Chain.GetStatus()
	if err != nil {
		return nil, err
	}

	return ToDetailsView(mostRecentSyncable, getHeadRes, getStatusRes), nil
}