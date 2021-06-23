package indexing

import (
	"context"

	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

type ReindexCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *reindexUseCase
}

func NewReindexCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *ReindexCmdHandler {
	return &ReindexCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *ReindexCmdHandler) Handle(ctx context.Context, parallel bool, force bool, startHeight, endHeight int64, targetIds []int64) {
	logger.Info("running reindex use case [handler=cmd]")

	useCaseConfig := ReindexUseCaseConfig{
		Parallel:    parallel,
		Force:       force,
		StartHeight: startHeight,
		EndHeight:   endHeight,
		TargetIds:   targetIds,
	}
	err := h.getUseCase().Execute(ctx, useCaseConfig)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *ReindexCmdHandler) getUseCase() *reindexUseCase {
	if h.useCase == nil {
		h.useCase = NewReindexUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}
