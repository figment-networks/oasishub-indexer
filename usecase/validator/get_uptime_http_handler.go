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
	_ types.HttpHandler = (*getUptimeHttpHandler)(nil)
)

type getUptimeHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getUptimeUseCase
}

func NewGetUptimeHttpHandler(db *store.Store, c *client.Client) *getUptimeHttpHandler {
	return &getUptimeHttpHandler{
		db:     db,
		client: c,
	}
}

type GetUptimeRequest struct {
	EntityUID string `form:"entity_uid" binding:"required"`
	Interval  string `form:"interval" binding:"required"`
	Period    string `form:"period" binding:"required"`
}

func (h *getUptimeHttpHandler) Handle(c *gin.Context) {
	var req GetUptimeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid request parameters")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.getUseCase().Execute(req.EntityUID, req.Interval, req.Period)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *getUptimeHttpHandler) getUseCase() *getUptimeUseCase {
	if h.useCase == nil {
		return NewGetUptimeUseCase(h.db)
	}
	return h.useCase
}
