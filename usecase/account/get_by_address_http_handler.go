package account

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/http"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

type Request struct {
	Address string `uri:"address" binding:"required"`
	Height  int64  `form:"height" binding:"-"`
}

func (h *getByAddressHttpHandler) Handle(c *gin.Context) {
	var req Request
	if err := c.ShouldBindUri(&req); err != nil {
		http.BadRequest(c, errors.New("invalid address"))
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, errors.New("invalid height"))
		return
	}

	resp, err := h.getUseCase().Execute(req.Address, req.Height)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getByAddressHttpHandler) getUseCase() *getByAddressUseCase {
	if h.useCase == nil {
		h.useCase = NewGetByAddressUseCase(h.db, h.client)
	}
	return h.useCase
}
