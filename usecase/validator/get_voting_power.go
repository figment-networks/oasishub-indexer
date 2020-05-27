package validator

import (
	"github.com/figment-networks/oasishub-indexer/store"
)

type getVotingPowerUseCase struct {
	db *store.Store
}

func NewGetVotingPowerUseCase(db *store.Store) *getVotingPowerUseCase {
	return &getVotingPowerUseCase{
		db: db,
	}
}

func (uc *getVotingPowerUseCase) Execute(key string, interval string, period string) ([]store.AvgForTimeIntervalRow, error) {
	return uc.db.ValidatorSeq.GetValidatorVotingPowerForInterval(key, interval, period)
}

