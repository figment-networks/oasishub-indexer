package indexing

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/indexer"
	"github.com/figment-networks/oasishub-indexer/store"
)

type reindexUseCase struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client
}

func NewReindexUseCase(cfg *config.Config, db *store.Store, c *client.Client) *reindexUseCase {
	return &reindexUseCase{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

type ReindexUseCaseConfig struct {
	Parallel    bool
	StartHeight int64
	EndHeight   int64
	TargetIds   []int64
}

func (uc *reindexUseCase) Execute(ctx context.Context, useCaseConfig ReindexUseCaseConfig) error {
	indexingPipeline, err := indexer.NewPipeline(uc.cfg, uc.db, uc.client)
	if err != nil {
		return err
	}

	return indexingPipeline.Reindex(ctx, indexer.ReindexConfig{
		Parallel:    useCaseConfig.Parallel,
		StartHeight: useCaseConfig.StartHeight,
		EndHeight:   useCaseConfig.EndHeight,
		TargetIds:   useCaseConfig.TargetIds,
	})
}
