package block

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getBlockTimesForIntervalUseCase struct {
	db *store.Store
}

func NewGetBlockTimeForIntervalUseCase(db *store.Store) *getBlockTimesForIntervalUseCase {
	return &getBlockTimesForIntervalUseCase{
		db: db,
	}
}

func (uc *getBlockTimesForIntervalUseCase) Execute(interval string, period string) ([]store.GetAvgTimesForIntervalRow, error) {
	return uc.db.BlockSeq.GetAvgTimesForInterval(interval, period)
}
