package validator

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getUptimeUseCase struct {
	db *store.Store
}

func NewGetUptimeUseCase(db *store.Store) *getUptimeUseCase {
	return &getUptimeUseCase{
		db: db,
	}
}

func (uc *getUptimeUseCase) Execute(key string, interval string, period string) ([]store.AvgForTimeIntervalRow, error) {
	return uc.db.ValidatorSeq.GetValidatorUptimeForInterval(key, interval, period)
}
