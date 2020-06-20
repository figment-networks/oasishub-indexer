package validator

import (
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
)

type getSummaryUseCase struct {
	db *store.Store
}

func NewGetSummaryUseCase(db *store.Store) *getSummaryUseCase {
	return &getSummaryUseCase{
		db: db,
	}
}

func (uc *getSummaryUseCase) Execute(interval types.SummaryInterval, period string, address string) (interface{}, error) {
	if address == "" {
		return uc.db.ValidatorSummary.FindSummary(interval, period)
	}
	return uc.db.ValidatorSummary.FindSummaryByAddress(address, interval, period)
}

