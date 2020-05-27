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
	_ types.HttpHandler = (*getVotingPowerHttpHandler)(nil)
)

type getVotingPowerHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getVotingPowerUseCase
}

func NewGetVotingPowerHttpHandler(db *store.Store, c *client.Client) *getVotingPowerHttpHandler {
	return &getVotingPowerHttpHandler{
		db:     db,
		client: c,
	}
}

type Request struct {
	EntityUID string `form:"entity_uid" binding:"required"`
	Interval  string `form:"interval" binding:"required"`
	Period    string `form:"period" binding:"required"`
}

func (h *getVotingPowerHttpHandler) Handle(c *gin.Context) {
	var req Request
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

func (h *getVotingPowerHttpHandler) getUseCase() *getVotingPowerUseCase {
	if h.useCase == nil {
		return NewGetVotingPowerUseCase(h.db)
	}
	return h.useCase
}
