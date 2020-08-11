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
		IndexerIndex:     indexing.NewIndexWorkerHandler(cfg, db, c),
		IndexerSummarize: indexing.NewSummarizeWorkerHandler(cfg, db, c),
		IndexerPurge:     indexing.NewPurgeWorkerHandler(cfg, db, c),
	}
}

type WorkerHandlers struct {
	IndexerIndex     types.WorkerHandler
	IndexerSummarize types.WorkerHandler
	IndexerPurge     types.WorkerHandler
}
