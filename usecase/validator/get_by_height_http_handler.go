package validator

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/config"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/http"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	_ types.HttpHandler = (*getByHeightHttpHandler)(nil)
)

type getByHeightHttpHandler struct {
	cfg    *config.Config
	db     *store.Store
	client *client.Client

	useCase *getByHeightUseCase
}

func NewGetByHeightHttpHandler(cfg *config.Config, db *store.Store, c *client.Client) *getByHeightHttpHandler {
	return &getByHeightHttpHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

type GetByHeightRequest struct {
	Height *int64 `form:"height" binding:"-"`
}

func (h *getByHeightHttpHandler) Handle(c *gin.Context) {
	var req GetByHeightRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, errors.New("invalid height"))
		return
	}

	resp, err := h.getUseCase().Execute(req.Height)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getByHeightHttpHandler) getUseCase() *getByHeightUseCase {
	if h.useCase == nil {
		h.useCase = NewGetByHeightUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}
