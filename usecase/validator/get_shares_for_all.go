package validator

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getSharesForAllUseCase struct {
	db *store.Store
}

func NewGetSharesForAllUseCase(db *store.Store) *getSharesForAllUseCase {
	return &getSharesForAllUseCase{
		db: db,
	}
}

func (uc *getSharesForAllUseCase) Execute(interval string, period string) ([]store.AvgQuantityForTimeIntervalRow, error) {
	return uc.db.ValidatorSeq.GetTotalSharesForInterval(interval, period)
}

