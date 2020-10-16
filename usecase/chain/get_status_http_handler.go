package chain

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/http"
	"github.com/gin-gonic/gin"
)

var (
	_ types.HttpHandler = (*getStatusHttpHandler)(nil)
)

type getStatusHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getStatusUseCase
}

func NewGetStatusHttpHandler(db *store.Store, client *client.Client) *getStatusHttpHandler {
	return &getStatusHttpHandler{
		db:     db,
		client: client,
	}
}

func (h *getStatusHttpHandler) Handle(c *gin.Context) {
	resp, err := h.getUseCase().Execute(c)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getStatusHttpHandler) getUseCase() *getStatusUseCase {
	if h.useCase == nil {
		h.useCase = NewGetStatusUseCase(h.db, h.client)
	}
	return h.useCase
}
