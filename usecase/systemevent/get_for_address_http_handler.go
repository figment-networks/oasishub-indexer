package systemevent

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

var (
	_ types.HttpHandler = (*getForAddressHttpHandler)(nil)
)

type getForAddressHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getForAddressUseCase
}

func NewGetForAddressHttpHandler(db *store.Store, c *client.Client) *getForAddressHttpHandler {
	return &getForAddressHttpHandler{
		db:     db,
		client: c,
	}
}

type GetForAddressRequest struct {
	Address string                 `uri:"address" binding:"required"`
	After   *int64                 `form:"after" binding:"-"`
	Kind    *model.SystemEventKind `form:"kind" binding:"-"`
}

func (h *getForAddressHttpHandler) Handle(c *gin.Context) {
	var req GetForAddressRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid address")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid kind or/and after")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.getUseCase().Execute(req.Address, req.After, req.Kind)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *getForAddressHttpHandler) getUseCase() *getForAddressUseCase {
	if h.useCase == nil {
		h.useCase = NewGetForAddressUseCase(h.db)
	}
	return h.useCase
}