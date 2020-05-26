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

func (uc *getByPublicKeyUseCase) Execute(key string) (*DetailsView, error) {
	rawAccount, err := uc.client.Account.GetByPublicKey(key)
	if err != nil {
		return nil, err
	}

	accountAgg, err := uc.db.AccountAgg.FindByPublicKey(key)
	if err != nil {
		return nil, err
	}

	delegationSeqs, err := uc.db.DelegationSeq.FindCurrentByDelegatorUID(key)
	if err != nil {
		return nil, err
	}

	debondingDelegationsSeqs, err := uc.db.DebondingDelegationSeq.FindRecentByDelegatorUID(key, 5)
	if err != nil {
		return nil, err
	}

	return ToDetailsView(rawAccount.GetAccount(), accountAgg, delegationSeqs, debondingDelegationsSeqs), nil
}
