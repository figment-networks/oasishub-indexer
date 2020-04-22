package getmostrecentheight

import (
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpHandler struct {
	useCase UseCase
}

type Response struct {
	Height types.Height `json:"height"`
}

func NewHttpHandler(useCase UseCase) types.HttpHandler {
	return &httpHandler{useCase: useCase}
}

func (h *httpHandler) Handle(c *gin.Context) {
	height, err := h.useCase.Execute()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, Response{Height: *height})
}


