package validator

import (
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

func (uc *getByEntityUidUseCase) Execute(key string) (*AggDetailsView, error) {
	validatorAggs, err := uc.db.ValidatorAgg.FindByEntityUID(key)
	if err != nil {
		return nil, err
	}

	return ToAggDetailsView(validatorAggs), nil
}

