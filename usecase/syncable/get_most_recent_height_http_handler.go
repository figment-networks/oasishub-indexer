package syncable

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	_ types.HttpHandler = (*getMostRecentHeightHttpHandler)(nil)
)

type getMostRecentHeightHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getMostRecentHeightUseCase
}

type Response struct {
	Height int64 `json:"height"`
}

func NewGetMostRecentHeightHttpHandler(db *store.Store, c *client.Client) *getMostRecentHeightHttpHandler {
	return &getMostRecentHeightHttpHandler{
		db: db,
		client: c,
	}
}

func (h *getMostRecentHeightHttpHandler) Handle(c *gin.Context) {
	height, err := h.getUseCase().Execute()
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, Response{Height: *height})
}

func (h *getMostRecentHeightHttpHandler) getUseCase() *getMostRecentHeightUseCase {
	if h.useCase == nil {
		return NewGetMostRecentHeightUseCase(h.db)
	}
	return h.useCase
}

