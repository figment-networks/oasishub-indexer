package getaccountbypublickey

import (
	"github.com/figment-networks/oasishub-indexer/mappers/accountaggmapper"
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
	PublicKey types.PublicKey `uri:"public_key" binding:"required"`
}

func (h *httpHandler) Handle(c *gin.Context) {
	var req Request
	if err := c.ShouldBindUri(&req); err != nil {
		log.Error(err)
		err := errors.NewError("invalid public key", errors.ServerInvalidParamsError, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.useCase.Execute(req.PublicKey)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, accountaggmapper.ToView(res))
}
