package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/indexing"
)

func NewWorkerHandlers(cfg *config.Config, db *store.Store, c *client.Client) *WorkerHandlers {
	return &WorkerHandlers{
		RunIndexer:       indexing.NewRunWorkerHandler(cfg, db, c),
		SummarizeIndexer: indexing.NewSummarizeWorkerHandler(cfg, db, c),
		PurgeIndexer:     indexing.NewPurgeWorkerHandler(cfg, db, c),
	}
}

type WorkerHandlers struct {
	RunIndexer       types.WorkerHandler
	SummarizeIndexer types.WorkerHandler
	PurgeIndexer     types.WorkerHandler
}
