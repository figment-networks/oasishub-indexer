package validator

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getSharesUseCase struct {
	db *store.Store
}

func NewGetSharesUseCase(db *store.Store) *getSharesUseCase {
	return &getSharesUseCase{
		db: db,
	}
}

func (uc *getSharesUseCase) Execute(key string, interval string, period string) ([]store.AvgQuantityForTimeIntervalRow, error) {
	return uc.db.ValidatorSeq.GetValidatorSharesForInterval(key, interval, period)
}

