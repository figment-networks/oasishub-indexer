package indexing

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/indexer"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/pkg/errors"
)

var (
	ErrBackfillRunning = errors.New("backfill already running (use force flag to override it)")
)

type backfillUseCase struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client
}

func NewBackfillUseCase(cfg *config.Config, db *store.Store, c *client.Client) *backfillUseCase {
	return &backfillUseCase{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

type BackfillUseCaseConfig struct {
	Parallel   bool
	Force      bool
}

func (uc *backfillUseCase) Execute(ctx context.Context, useCaseConfig BackfillUseCaseConfig) error {
	indexingPipeline := indexer.NewPipeline(uc.cfg, uc.db, uc.client)

	return indexingPipeline.Backfill(ctx, indexer.BackfillConfig{
		Parallel:   useCaseConfig.Parallel,
		Force:      useCaseConfig.Force,
	})
}