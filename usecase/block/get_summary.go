package block

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getBlockSummaryUseCase struct {
	db *store.Store
}

func NewGetBlockSummaryUseCase(db *store.Store) *getBlockSummaryUseCase {
	return &getBlockSummaryUseCase{
		db: db,
	}
}

func (uc *getBlockSummaryUseCase) Execute(interval string, period string) ([]store.GetAvgTimesForIntervalRow, error) {
	return uc.db.BlockSeq.GetSummary(interval, period)
}
