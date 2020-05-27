package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

var (
	_ types.WorkerHandler = (*runIndexerWorkerHandler)(nil)
)

type runIndexerWorkerHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *runIndexerUseCase
}

func NewRunIndexerWorkerHandler(cfg *config.Config, db *store.Store, c *client.Client) *runIndexerWorkerHandler {
	return &runIndexerWorkerHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *runIndexerWorkerHandler) Handle() {
	batchSize := h.cfg.DefaultBatchSize
	ctx := context.Background()

	logger.Info(fmt.Sprintf("running indexer use case [handler=worker] [batchSize=%d]", batchSize))

	err := h.getUseCase().Execute(ctx, batchSize)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *runIndexerWorkerHandler) getUseCase() *runIndexerUseCase {
	if h.useCase == nil {
		return NewRunIndexerUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}


