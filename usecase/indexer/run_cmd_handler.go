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
	_ types.CmdHandler = (*runCmdHandler)(nil)
)

type runCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *runUseCase
}

func NewRunCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *runCmdHandler {
	return &runCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *runCmdHandler) Handle(ctx context.Context) {
	//TODO: Pass as an argument from command line
	batchSize := int64(8000)

	logger.Info(fmt.Sprintf("running indexer use case [handler=cmd] [batchSize=%d]", batchSize))

	err := h.getUseCase().Execute(ctx, batchSize)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *runCmdHandler) getUseCase() *runUseCase {
	if h.useCase == nil {
		return NewRunUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}
