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
	_ types.CmdHandler = (*summarizeCmdHandler)(nil)
)

type summarizeCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *summarizeUseCase
}

func NewSummarizeCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *summarizeCmdHandler {
	return &summarizeCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *summarizeCmdHandler) Handle(ctx context.Context) {
	logger.Info(fmt.Sprintf("summarizing indexer use case [handler=cmd]"))

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *summarizeCmdHandler) getUseCase() *summarizeUseCase {
	if h.useCase == nil {
		return NewSummarizeUseCase(h.cfg, h.db)
	}
	return h.useCase
}

