package debondingdelegation

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/http"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

type Request struct {
	Height *int64 `form:"height" binding:"-"`
}

func (h *getByHeightHttpHandler) Handle(c *gin.Context) {
	var req Request
	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, errors.New("invalid height"))
		return
	}

	resp, err := h.getUseCase().Execute(req.Height)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getByHeightHttpHandler) getUseCase() *getByHeightUseCase {
	if h.useCase == nil {
		h.useCase = NewGetByHeightUseCase(h.db, h.client)
	}
	return h.useCase
}



