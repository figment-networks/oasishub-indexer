package balance

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
)

type getForAddressUseCase struct {
	db *store.Store
}

func NewGetForAddressUseCase(db *store.Store) *getForAddressUseCase {
	return &getForAddressUseCase{
		db: db,
	}
}

func (uc *getForAddressUseCase) Execute(address string, start, end *types.Time) ([]model.BalanceSummary, error) {
	summaries, err := uc.db.BalanceSummary.GetDailySummaries(address, start, end)
	if err != nil {
		return nil, err
	}

	if len(summaries) == 0 {
		_, err := uc.db.ValidatorAgg.FindByAddress(address)
		return summaries, err
	}

	return summaries, nil
}
