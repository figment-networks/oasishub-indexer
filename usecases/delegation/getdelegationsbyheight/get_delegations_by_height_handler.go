package getdelegationsbyheight

import (
	"github.com/figment-networks/oasishub-indexer/mappers/delegationseqmapper"
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
	Height *types.Height `form:"height" binding:"-"`
}

func (h *httpHandler) Handle(c *gin.Context) {
	var req Request
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Error(err)
		err := errors.NewError("invalid height", errors.ServerInvalidParamsError, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ds, err := h.useCase.Execute(req.Height)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, delegationseqmapper.ToView(ds))
}


