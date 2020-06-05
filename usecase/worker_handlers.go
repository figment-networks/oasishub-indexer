package usecase

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/indexer"
)

func NewWorkerHandlers(cfg *config.Config, db *store.Store, c *client.Client) *WorkerHandlers {
	return &WorkerHandlers{
		RunIndexer:   indexer.NewRunWorkerHandler(cfg, db, c),
		PurgeIndexer: indexer.NewPurgeWorkerHandler(cfg, db, c),
	}
}

type WorkerHandlers struct {
	RunIndexer   types.WorkerHandler
	PurgeIndexer types.WorkerHandler
}
