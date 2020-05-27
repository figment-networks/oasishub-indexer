package block

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
	_ types.HttpHandler = (*getBlockTimesForIntervalHttpHandler)(nil)
)

type getBlockTimesForIntervalHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getBlockTimesForIntervalUseCase
}

func NewGetBlockTimesForIntervalHttpHandler(db *store.Store, client *client.Client) *getBlockTimesForIntervalHttpHandler {
	return &getBlockTimesForIntervalHttpHandler{
		db:     db,
		client: client,
	}
}

type GetBlockTimesForIntervalRequest struct {
	Interval string `form:"interval" binding:"required"`
	Period   string `form:"period" binding:"required"`
}

func (h *getBlockTimesForIntervalHttpHandler) Handle(c *gin.Context) {
	var req GetBlockTimesForIntervalRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid interval and/or period")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.getUseCase().Execute(req.Interval, req.Period)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *getBlockTimesForIntervalHttpHandler) getUseCase() *getBlockTimesForIntervalUseCase {
	if h.useCase == nil {
		return NewGetBlockTimeForIntervalUseCase(h.db)
	}
	return h.useCase
}