package validator

import (
	"context"

	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

type DecorateCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *decorateUseCase
}

func NewDecorateCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *DecorateCmdHandler {
	return &DecorateCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *DecorateCmdHandler) Handle(ctx context.Context, filePath string) {
	logger.Info("running decorate validators indexer use case [handler=cmd]")

	err := h.getUseCase().Execute(ctx, filePath)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *DecorateCmdHandler) getUseCase() *decorateUseCase {
	if h.useCase == nil {
		return NewDecorateUseCase(h.cfg, h.db.ValidatorAgg)
	}
	return h.useCase
}
