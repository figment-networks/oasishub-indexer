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
	_ types.CmdHandler = (*parseCSVCmdHandler)(nil)
)

type parseCSVCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *parseCSVUseCase
}

func NewParseCSVCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *parseCSVCmdHandler {
	return &parseCSVCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *parseCSVCmdHandler) Handle(ctx context.Context) {
	logger.Info(fmt.Sprintf("running parse csv indexer use case [handler=cmd]"))

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *parseCSVCmdHandler) getUseCase() *parseCSVUseCase {
	if h.useCase == nil {
		return NewParseCSVUseCase(h.cfg, h.db)
	}
	return h.useCase
}
