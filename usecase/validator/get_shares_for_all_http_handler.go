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
	_ types.HttpHandler = (*getSharesForAllHttpHandler)(nil)
)

type getSharesForAllHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getSharesForAllUseCase
}

func NewGetSharesForAllHttpHandler(db *store.Store, c *client.Client) *getSharesForAllHttpHandler {
	return &getSharesForAllHttpHandler{
		db:     db,
		client: c,
	}
}

type GetSharesForAllRequest struct {
	Interval string `form:"interval" binding:"required"`
	Period   string `form:"period" binding:"required"`
}

func (h *getSharesForAllHttpHandler) Handle(c *gin.Context) {
	var req GetSharesForAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid request parameters")
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

func (h *getSharesForAllHttpHandler) getUseCase() *getSharesForAllUseCase {
	if h.useCase == nil {
		return NewGetSharesForAllUseCase(h.db)
	}
	return h.useCase
}

