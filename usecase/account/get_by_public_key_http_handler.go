package account

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
	_ types.HttpHandler = (*getByPublicKeyHttpHandler)(nil)
)

type getByPublicKeyHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getByPublicKeyUseCase
}

func NewGetByPublicKeyHttpHandler(db *store.Store, c *client.Client) *getByPublicKeyHttpHandler {
	return &getByPublicKeyHttpHandler{
		db:     db,
		client: c,
	}
}

type Request struct {
	PublicKey string `form:"public_key" binding:"required"`
	Height    int64  `form:"height" binding:"-"`
}

func (h *getByPublicKeyHttpHandler) Handle(c *gin.Context) {
	var req Request
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(err)
		err := errors.New("invalid public key")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.getUseCase().Execute(req.PublicKey, req.Height)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *getByPublicKeyHttpHandler) getUseCase() *getByPublicKeyUseCase {
	if h.useCase == nil {
		return NewGetByPublicKeyUseCase(h.db, h.client)
	}
	return h.useCase
}
