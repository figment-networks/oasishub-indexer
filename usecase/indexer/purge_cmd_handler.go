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
	_ types.CmdHandler = (*purgeCmdHandler)(nil)
)

type purgeCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *purgeUseCase
}

func NewPurgeCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *purgeCmdHandler {
	return &purgeCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *purgeCmdHandler) Handle(ctx context.Context) {
	logger.Info("running indexer use case [handler=cmd]")

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *purgeCmdHandler) getUseCase() *purgeUseCase {
	if h.useCase == nil {
		return NewPurgeUseCase(h.cfg, h.db)
	}
	return h.useCase
}
