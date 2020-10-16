package transaction

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
	_ types.HttpHandler = (*broadcastHttpHandler)(nil)
)

type broadcastHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *broadcastUseCase
}

func NewBroadcastHttpHandler(db *store.Store, c *client.Client) *broadcastHttpHandler {
	return &broadcastHttpHandler{
		db: db,
		client: c,
	}
}

type BroadcastRequest struct {
	TxRaw string `form:"tx_raw" binding:"required" json:"tx_raw"`
}

type BroadcastResponse struct {
	Submitted bool `json:"submitted"`
}

func (h *broadcastHttpHandler) Handle(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid raw transaction string")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.getUseCase().Execute(req.TxRaw)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, BroadcastResponse{Submitted: *res})
}

func (h *broadcastHttpHandler) getUseCase() *broadcastUseCase {
	if h.useCase == nil {
		h.useCase = NewBroadcastUseCase(h.db, h.client)
	}
	return h.useCase
}



