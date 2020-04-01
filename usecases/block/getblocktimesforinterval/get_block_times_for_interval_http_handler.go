package getblocktimesforinterval

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	InvalidParamsError = "block_handler_invalid_params_error"
)

type httpHandler struct {
	useCase UseCase
}

func NewHttpHandler(useCase UseCase) types.HttpHandler {
	return &httpHandler{useCase: useCase}
}

type Request struct {
	Interval string `uri:"interval" binding:"required"`
}

func (h *httpHandler) Handle(c *gin.Context) {
	var req Request
	if err := c.ShouldBindUri(&req); err != nil {
		log.Error(err)
		err := errors.NewError("invalid interval", errors.ServerInvalidParamsError, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.useCase.Execute(req.Interval)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

