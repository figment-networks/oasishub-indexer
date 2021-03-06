package block

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getBlockTimesUseCase struct {
	db *store.Store
}

func NewGetBlockTimesUseCase(db *store.Store) *getBlockTimesUseCase {
	return &getBlockTimesUseCase{
		db: db,
	}
}

func (uc *getBlockTimesUseCase) Execute(limit int64) (*store.GetAvgRecentTimesResult, error) {
	return uc.db.BlockSeq.GetAvgRecentTimes(limit)
}

