package indexing

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
	_ types.WorkerHandler = (*indexWorkerHandler)(nil)
)

type indexWorkerHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *indexUseCase
}

func NewIndexWorkerHandler(cfg *config.Config, db *store.Store, c *client.Client) *indexWorkerHandler {
	return &indexWorkerHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *indexWorkerHandler) Handle() {
	batchSize := h.cfg.DefaultBatchSize
	ctx := context.Background()

	logger.Info(fmt.Sprintf("running indexer use case [handler=worker] [batchSize=%d]", batchSize))

	err := h.getUseCase().Execute(ctx, batchSize)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *indexWorkerHandler) getUseCase() *indexUseCase {
	if h.useCase == nil {
		h.useCase = NewIndexUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}


