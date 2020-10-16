package validator

import (
	"github.com/figment-networks/oasishub-indexer/client"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/usecase/http"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	_ types.HttpHandler = (*getSummaryHttpHandler)(nil)

	ErrInvalidIntervalPeriod = errors.New("invalid interval and/or period")
)

type getSummaryHttpHandler struct {
	db     *store.Store
	client *client.Client

	useCase *getSummaryUseCase
}

func NewGetSummaryHttpHandler(db *store.Store, c *client.Client) *getSummaryHttpHandler {
	return &getSummaryHttpHandler{
		db:     db,
		client: c,
	}
}

type GetSummaryRequest struct {
	Interval types.SummaryInterval `form:"interval" binding:"required"`
	Period   string                `form:"period" binding:"required"`
	Address  string                `form:"address" binding:"-"`
}

func (h *getSummaryHttpHandler) Handle(c *gin.Context) {
	req, err := h.validateParams(c)
	if err != nil {
		http.BadRequest(c, err)
		return
	}

	resp, err := h.getUseCase().Execute(req.Interval, req.Period, req.Address)
	if http.ShouldReturn(c, err) {
		return
	}

	http.JsonOK(c, resp)
}

func (h *getSummaryHttpHandler) validateParams(c *gin.Context) (*GetSummaryRequest, error) {
	var req GetSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	if !req.Interval.Valid() {
		return nil, ErrInvalidIntervalPeriod
	}

	return &req, nil
}

func (h *getSummaryHttpHandler) getUseCase() *getSummaryUseCase {
	if h.useCase == nil {
		h.useCase = NewGetSummaryUseCase(h.db)
	}
	return h.useCase
}
