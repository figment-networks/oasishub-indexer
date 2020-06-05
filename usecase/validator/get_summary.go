package validator

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getSummaryUseCase struct {
	db *store.Store
}

func NewGetSummaryUseCase(db *store.Store) *getSummaryUseCase {
	return &getSummaryUseCase{
		db: db,
	}
}

func (uc *getSummaryUseCase) Execute(interval string, period string, entityUID string) (interface{}, error) {
	if entityUID == "" {
		return uc.db.ValidatorSeq.GetSummary(interval, period)
	}
	return uc.db.ValidatorSeq.GetSummaryByEntityUID(entityUID, interval, period)
}

