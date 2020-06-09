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

func (uc *getSummaryUseCase) Execute(interval types.SummaryInterval, period string, entityUID string) (interface{}, error) {
	if entityUID == "" {
		return uc.db.ValidatorSummary.FindSummary(interval, period)
	}
	return uc.db.ValidatorSummary.FindSummaryByEntityUID(entityUID, interval, period)
}

