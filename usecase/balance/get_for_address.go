package balance

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
)

type getForAddressUseCase struct {
	db *store.Store
}

func NewGetForAddressUseCase(db *store.Store) *getForAddressUseCase {
	return &getForAddressUseCase{
		db: db,
	}
}

func (uc *getForAddressUseCase) Execute(address, start, end string) ([]model.BalanceSummary, error) {
	summaries, err := uc.db.BalanceSummary.GetDailySummaries(address, start, end)
	if err != nil {
		return nil, err
	}

	return summaries, nil
}
