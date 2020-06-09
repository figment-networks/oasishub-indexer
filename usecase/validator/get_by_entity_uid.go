package validator

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
)

type getByEntityUidUseCase struct {
	db *store.Store
}

func NewGetByEntityUidUseCase(db *store.Store) *getByEntityUidUseCase {
	return &getByEntityUidUseCase{
		db: db,
	}
}

func (uc *getByEntityUidUseCase) Execute(key string, sequencesLimit int64) (*AggDetailsView, error) {
	validatorAggs, err := uc.db.ValidatorAgg.FindByEntityUID(key)
	if err != nil {
		return nil, err
	}

	sequences, err := uc.getSequences(key, sequencesLimit)
	if err != nil {
		return nil, err
	}

	return ToAggDetailsView(validatorAggs, sequences), nil
}

func (uc *getByEntityUidUseCase) getSequences(key string, sequencesLimit int64) ([]model.ValidatorSeq, error) {
	var sequences []model.ValidatorSeq
	var err error
	if sequencesLimit > 0 {
		sequences, err = uc.db.ValidatorSeq.FindLastByEntityUID(key, sequencesLimit)
		if err != nil {
			return nil, err
		}
	}
	return sequences, nil
}

