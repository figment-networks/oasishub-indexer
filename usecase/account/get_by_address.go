package account

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
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

func (uc *getByAddressUseCase) Execute(address string, height int64) (*DetailsView, error) {
	rawAccount, err := uc.client.Account.GetByAddress(address, height)
	if err != nil {
		return nil, err
	}

	return ToDetailsView(rawAccount.GetAccount()), nil
}
