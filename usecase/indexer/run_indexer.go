package indexer

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	indexing "github.com/figment-networks/oasishub-indexer/indexer"
	"github.com/figment-networks/oasishub-indexer/store"
)

type runIndexerUseCase struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client
}

func NewRunIndexerUseCase(cfg *config.Config, db *store.Store, c *client.Client) *runIndexerUseCase {
	return &runIndexerUseCase{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (uc *runIndexerUseCase) Execute(ctx context.Context, batchSize int64) error {
	pipeline, err := indexing.NewPipeline(uc.cfg, uc.db, uc.client)
	if err != nil {
		return err
	}
	return pipeline.Start(ctx, batchSize)
}
