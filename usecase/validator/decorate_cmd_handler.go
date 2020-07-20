package validator

import (
	"context"
	"fmt"

	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

const (
	// CtxFilePath is context key for value containing file path
	CtxFilePath = "context_file"
)

var (
	_ types.CmdHandler = (*decorateCmdHandler)(nil)

	errMissingFileInCtx = errors.New("missing file path in context")
)

type decorateCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *decorateUseCase
}

func NewDecorateCmdHandler(cfg *config.Config, db *store.Store, c *client.Client) *decorateCmdHandler {
	return &decorateCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *decorateCmdHandler) Handle(ctx context.Context) {
	logger.Info(fmt.Sprintf("running decorate validators indexer use case [handler=cmd]"))

	filePath, ok := ctx.Value(CtxFilePath).(string)
	if !ok {
		logger.Error(errMissingFileInCtx)
		return
	}

	err := h.getUseCase().Execute(ctx, filePath)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *decorateCmdHandler) getUseCase() *decorateUseCase {
	if h.useCase == nil {
		return NewDecorateUseCase(h.cfg, h.db)
	}
	return h.useCase
}
