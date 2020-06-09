package validator

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

var (
	_ types.HttpHandler = (*getByEntityUidHttpHandler)(nil)
)

type getByEntityUidHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getByEntityUidUseCase
}

func NewGetByEntityUidHttpHandler(db *store.Store, c *client.Client) *getByEntityUidHttpHandler {
	return &getByEntityUidHttpHandler{
		db:     db,
		client: c,
	}
}

type GetByEntityUidRequest struct {
	EntityUID      string `form:"entity_uid" binding:"required"`
	SequencesLimit int64  `form:"sequences_limit" binding:"-"`
}

func (h *getByEntityUidHttpHandler) Handle(c *gin.Context) {
	var req GetByEntityUidRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid entity_uid")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.getUseCase().Execute(req.EntityUID, req.SequencesLimit)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *getByEntityUidHttpHandler) getUseCase() *getByEntityUidUseCase {
	if h.useCase == nil {
		return NewGetByEntityUidUseCase(h.db)
	}
	return h.useCase
}
