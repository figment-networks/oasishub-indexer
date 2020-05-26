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
	_ types.HttpHandler = (*getSharesHttpHandler)(nil)
)

type getSharesHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getSharesUseCase
}

func NewGetSharesHttpHandler(db *store.Store, c *client.Client) *getSharesHttpHandler {
	return &getSharesHttpHandler{
		db:     db,
		client: c,
	}
}

type GetSharesRequest struct {
	EntityUID string `form:"entity_uid" binding:"required"`
	Interval  string `form:"interval" binding:"required"`
	Period    string `form:"period" binding:"required"`
}

func (h *getSharesHttpHandler) Handle(c *gin.Context) {
	var req GetSharesRequest
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

func (h *getSharesHttpHandler) getUseCase() *getSharesUseCase {
	if h.useCase == nil {
		return NewGetSharesUseCase(h.db)
	}
	return h.useCase
}
