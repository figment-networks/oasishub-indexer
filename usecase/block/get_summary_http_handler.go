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
	_ types.HttpHandler = (*getBlockSummaryHttpHandler)(nil)

	ErrInvalidIntervalPeriod = errors.New("invalid interval and/or period")
)

type getBlockSummaryHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getBlockSummaryUseCase
}

func NewGetBlockSummaryHttpHandler(db *store.Store, client *client.Client) *getBlockSummaryHttpHandler {
	return &getBlockSummaryHttpHandler{
		db:     db,
		client: client,
	}
}

type GetBlockTimesForIntervalRequest struct {
	Interval string `form:"interval" binding:"required"`
	Period   string `form:"period" binding:"required"`
}

func (h *getBlockSummaryHttpHandler) Handle(c *gin.Context) {
	req, err := h.validateParams(c)
	if err != nil {
		logger.Error(err)
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

func (h *getBlockSummaryHttpHandler) validateParams(c *gin.Context) (*GetBlockTimesForIntervalRequest, error) {
	var req GetBlockTimesForIntervalRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	if req.Interval != "hourly" && req.Interval != "daily" {
		return nil, ErrInvalidIntervalPeriod
	}

	return &req, nil
}

func (h *getBlockSummaryHttpHandler) getUseCase() *getBlockSummaryUseCase {
	if h.useCase == nil {
		return NewGetBlockSummaryUseCase(h.db)
	}
	return h.useCase
}
