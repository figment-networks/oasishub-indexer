package account

import (
	"time"

	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/http"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	_ types.HttpHandler = (*getSummariesHttpHandler)(nil)
)

type getSummariesHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getSummariesUseCase
}

func NewGetSummariesHttpHandler(db *store.Store, c *client.Client) *getSummariesHttpHandler {
	return &getSummariesHttpHandler{
		db:     db,
		client: c,
	}
}

type uriParams struct {
	Address string `uri:"address" binding:"required"`
}
type queryParams struct {
	Start time.Time `form:"start" binding:"required" time_format:"2006-01-02 15:04:05"`
	End   time.Time `form:"end" binding:"required" time_format:"2006-01-02 15:04:05"`
}

func (h *getSummariesHttpHandler) Handle(c *gin.Context) {
	var uri uriParams
	if err := c.ShouldBindUri(&uri); err != nil {
		http.BadRequest(c, errors.New("invalid address"))
		return
	}

	var params queryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		http.BadRequest(c, errors.New("invalid start and/or end params: must be in format \"2006-01-02 15:04:05\""))
		return
	}

	resp, err := h.getUseCase().Execute(uri.Address, params.Start, params.End)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getSummariesHttpHandler) getUseCase() *getSummariesUseCase {
	if h.useCase == nil {
		h.useCase = NewGetSummariesUseCase(h.db, h.client)
	}
	return h.useCase
}
