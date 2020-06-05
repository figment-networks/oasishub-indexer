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
	_ types.WorkerHandler = (*runWorkerHandler)(nil)
)

type runWorkerHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *runUseCase
}

func NewRunWorkerHandler(cfg *config.Config, db *store.Store, c *client.Client) *runWorkerHandler {
	return &runWorkerHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *runWorkerHandler) Handle() {
	batchSize := h.cfg.DefaultBatchSize
	ctx := context.Background()

	logger.Info(fmt.Sprintf("running indexer use case [handler=worker] [batchSize=%d]", batchSize))

	err := h.getUseCase().Execute(ctx, batchSize)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *runWorkerHandler) getUseCase() *runUseCase {
	if h.useCase == nil {
		return NewRunUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}


