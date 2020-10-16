package balance

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
	Address string    `uri:"address" binding:"required"`
	Start   time.Time `form:"start" binding:"-" time_format:"2006-01-02"`
	End     time.Time `form:"end" binding:"-" time_format:"2006-01-02"`
}

func (h *getForAddressHttpHandler) Handle(c *gin.Context) {
	var req GetForAddressRequest
	if err := c.ShouldBindUri(&req); err != nil {
		http.BadRequest(c, errors.New("invalid address"))
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		http.BadRequest(c, errors.New("invalid start or/and end"))
		return
	}

	resp, err := h.getUseCase().Execute(req.Address, types.NewTimeFromTime(req.Start), types.NewTimeFromTime(req.End))
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getForAddressHttpHandler) getUseCase() *getForAddressUseCase {
	if h.useCase == nil {
		h.useCase = NewGetForAddressUseCase(h.db)
	}
	return h.useCase
}
