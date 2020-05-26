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
	_ types.HttpHandler = (*getVotingPowerForAllHttpHandler)(nil)
)

type getVotingPowerForAllHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getVotingPowerForAllUseCase
}

func NewGetVotingPowerForAllHttpHandler(db *store.Store, c *client.Client) *getVotingPowerForAllHttpHandler {
	return &getVotingPowerForAllHttpHandler{
		db:     db,
		client: c,
	}
}

type GetVotingPowerForAllRequest struct {
	Interval string `form:"interval" binding:"required"`
	Period   string `form:"period" binding:"required"`
}

func (h *getVotingPowerForAllHttpHandler) Handle(c *gin.Context) {
	var req GetVotingPowerForAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid request parameters")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.useCase.Execute(req.Interval, req.Period)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *getVotingPowerForAllHttpHandler) getUseCase() *getVotingPowerForAllUseCase {
	if h.useCase == nil {
		return NewGetVotingPowerForAllUseCase(h.db)
	}
	return h.useCase
}


