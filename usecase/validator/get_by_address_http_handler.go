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
	_ types.HttpHandler = (*getByAddressHttpHandler)(nil)
)

type getByAddressHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getByAddressUseCase
}

func NewGetByAddressHttpHandler(db *store.Store, c *client.Client) *getByAddressHttpHandler {
	return &getByAddressHttpHandler{
		db:     db,
		client: c,
	}
}

type GetByEntityUidRequest struct {
	Address        string `uri:"address" binding:"required"`
	SequencesLimit int64  `form:"sequences_limit" binding:"-"`
}

func (h *getByAddressHttpHandler) Handle(c *gin.Context) {
	var req GetByEntityUidRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid address")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid sequences limit")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.getUseCase().Execute(req.Address, req.SequencesLimit)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *getByAddressHttpHandler) getUseCase() *getByAddressUseCase {
	if h.useCase == nil {
		h.useCase = NewGetByAddressUseCase(h.db)
	}
	return h.useCase
}
