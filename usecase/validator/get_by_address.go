package validator

import (
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
)

type getByAddressUseCase struct {
	db *store.Store
}

func NewGetByAddressUseCase(db *store.Store) *getByAddressUseCase {
	return &getByAddressUseCase{
		db: db,
	}
}

func (uc *getByAddressUseCase) Execute(key string, sequencesLimit int64) (*AggDetailsView, error) {
	validatorAggs, err := uc.db.ValidatorAgg.FindByAddress(key)
	if err != nil {
		return nil, err
	}

	sequences, err := uc.getSequences(key, sequencesLimit)
	if err != nil {
		return nil, err
	}

	return ToAggDetailsView(validatorAggs, sequences), nil
}

func (uc *getByAddressUseCase) getSequences(address string, sequencesLimit int64) ([]model.ValidatorSeq, error) {
	var sequences []model.ValidatorSeq
	var err error
	if sequencesLimit > 0 {
		sequences, err = uc.db.ValidatorSeq.FindLastByAddress(address, sequencesLimit)
		if err != nil {
			return nil, err
		}
	}
	return sequences, nil
}

