package validator

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getVotingPowerForAllUseCase struct {
	db *store.Store
}

func NewGetVotingPowerForAllUseCase(db *store.Store) *getVotingPowerForAllUseCase {
	return &getVotingPowerForAllUseCase{
		db: db,
	}
}

func (uc *getVotingPowerForAllUseCase) Execute(interval string, period string) ([]store.AvgForTimeIntervalRow, error) {
	return uc.db.ValidatorSeq.GetTotalVotingPowerForInterval(interval, period)
}
