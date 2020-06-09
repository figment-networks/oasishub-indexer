package block

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
)

type getBlockSummaryUseCase struct {
	db *store.Store
}

func NewGetBlockSummaryUseCase(db *store.Store) *getBlockSummaryUseCase {
	return &getBlockSummaryUseCase{
		db: db,
	}
}

func (uc *getBlockSummaryUseCase) Execute(interval types.SummaryInterval, period string) ([]model.BlockSummary, error) {
	return uc.db.BlockSummary.FindSummary(interval, period)
}
