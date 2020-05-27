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
	_ types.HttpHandler = (*getByHeightHttpHandler)(nil)
)

type getByHeightHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getByHeightUseCase
}

func NewGetByHeightHttpHandler(db *store.Store, c *client.Client) *getByHeightHttpHandler {
	return &getByHeightHttpHandler{
		db: db,
		client: c,
	}
}

type GetByHeightRequest struct {
	Height *int64 `form:"height" binding:"-"`
}

func (h *getByHeightHttpHandler) Handle(c *gin.Context) {
	var req GetByHeightRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid height")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ds, err := h.getUseCase().Execute(req.Height)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ds)
}

func (h *getByHeightHttpHandler) getUseCase() *getByHeightUseCase {
	if h.useCase == nil {
		return NewGetByHeightUseCase(h.db)
	}
	return h.useCase
}



