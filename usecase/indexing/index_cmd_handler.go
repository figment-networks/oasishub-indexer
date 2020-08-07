package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

type IndexCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *indexUseCase
}

func NewIndexCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *IndexCmdHandler {
	return &IndexCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *IndexCmdHandler) Handle(ctx context.Context, batchSize int64) {
	logger.Info(fmt.Sprintf("running indexer use case [handler=cmd] [batchSize=%d]", batchSize))

	err := h.getUseCase().Execute(ctx, batchSize)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *IndexCmdHandler) getUseCase() *indexUseCase {
	if h.useCase == nil {
		h.useCase = NewIndexUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}
