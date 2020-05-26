package indexer

import (
	"context"
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

var (
	_ types.CmdHandler = (*runIndexerCmdHandler)(nil)
)

type runIndexerCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *runIndexerUseCase
}

func NewRunIndexerCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *runIndexerCmdHandler {
	return &runIndexerCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *runIndexerCmdHandler) Handle(ctx context.Context) {
	//TODO: Pass as an argument from command line
	batchSize := h.cfg.DefaultBatchSize

	err := h.getUseCase().Execute(ctx, batchSize)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *runIndexerCmdHandler) getUseCase() *runIndexerUseCase {
	if h.useCase == nil {
		return NewRunIndexerUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}
