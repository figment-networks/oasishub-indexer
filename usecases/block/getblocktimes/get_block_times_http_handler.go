package getblocktimes

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpHandler struct {
	useCase UseCase
}

func NewHttpHandler(useCase UseCase) types.HttpHandler {
	return &httpHandler{useCase: useCase}
}

type Request struct {
	Limit int64 `uri:"limit" binding:"required"`
}

func (h *httpHandler) Handle(c *gin.Context) {
	var req Request
	if err := c.ShouldBindUri(&req); err != nil {
		log.Error(err)
		err := errors.NewError("invalid height", errors.ServerInvalidParamsError, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp := h.useCase.Execute(req.Limit)

	c.JSON(http.StatusOK, resp)
}

