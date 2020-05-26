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
	ea, err := uc.db.ValidatorAgg.FindByEntityUID(key)
	if err != nil {
		return nil, err
	}

	ds, err := uc.db.DelegationSeq.FindLastByValidatorUID(key)
	if err != nil {
		return nil, err
	}

	dds, err := uc.db.DebondingDelegationSeq.FindRecentByValidatorUID(key, 5)
	if err != nil {
		return nil, err
	}

	return ToAggDetailsView(*ea, ds, dds), nil
}

