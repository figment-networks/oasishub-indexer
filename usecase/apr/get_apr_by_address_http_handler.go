package apr

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
	_ types.HttpHandler = (*getAprByAddressHttpHandler)(nil)
)

type getAprByAddressHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getAprByAddressUseCase
}

func NewGetAprByAddressHttpHandler(db *store.Store, c *client.Client) *getAprByAddressHttpHandler {
	return &getAprByAddressHttpHandler{
		db:     db,
		client: c,
	}
}

type GetAprByAddressRequest struct {
	Address        string    `uri:"address" binding:"required"`
	IncludeDailies bool      `uri:"including_dailies" binding:"required"`
	Start          time.Time `form:"start" binding:"required" time_format:"2006-01-02"`
	End            time.Time `form:"end" binding:"required" time_format:"2006-01-02"`
}

func (h *getAprByAddressHttpHandler) Handle(c *gin.Context) {
	var req GetAprByAddressRequest
	if err := c.ShouldBindUri(&req); err != nil {
		http.BadRequest(c, errors.New("missing parameter"))
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, errors.New("invalid start or/and end"))
		return
	}

	resp, err := h.getUseCase().Execute(req.Address, types.NewTimeFromTime(req.Start), types.NewTimeFromTime(req.End), req.IncludeDailies)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getAprByAddressHttpHandler) getUseCase() *getAprByAddressUseCase {
	if h.useCase == nil {
		h.useCase = NewGetAprByAddressUseCase(h.db, h.client)
	}
	return h.useCase
}
