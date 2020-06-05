package account

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
)

type getByPublicKeyUseCase struct {
	db     *store.Store
	client *client.Client
}

func NewGetByPublicKeyUseCase(db *store.Store, c *client.Client) *getByPublicKeyUseCase {
	return &getByPublicKeyUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getByPublicKeyUseCase) Execute(key string, height int64) (*DetailsView, error) {
	rawAccount, err := uc.client.Account.GetByPublicKey(key, height)
	if err != nil {
		return nil, err
	}

	return ToDetailsView(rawAccount.GetAccount()), nil
}
