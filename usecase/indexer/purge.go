package indexer

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
)

type purgeUseCase struct {
	cfg    *config.Config
	db     *store.Store
}

func NewPurgeUseCase(cfg *config.Config, db *store.Store) *purgeUseCase {
	return &purgeUseCase{
		cfg:    cfg,
		db:     db,
	}
}

func (uc *purgeUseCase) Execute(ctx context.Context) error {
	if err := uc.db.BlockSeq.PurgeOldRecords(uc.cfg); err != nil {
		return err
	}

	if err := uc.db.ValidatorSeq.PurgeOldRecords(uc.cfg); err != nil {
		return err
	}

	return nil
}
