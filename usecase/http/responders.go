package http

import (
	"errors"
	"net/http"

	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"

	"github.com/gin-gonic/gin"
)

// BadRequest renders a HTTP 400 bad request response
func BadRequest(c *gin.Context, err error) {
	jsonError(c, http.StatusBadRequest, err)
}

// NotFound renders a HTTP 404 not found response
func NotFound(c *gin.Context, err error) {
	jsonError(c, http.StatusNotFound, err)
}

// ServerError renders a HTTP 500 error response
func ServerError(c *gin.Context, err error) {
	jsonError(c, http.StatusInternalServerError, err)
}

// JsonOK renders a successful response
func JsonOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// jsonError renders an error response
func jsonError(c *gin.Context, status int, err error) {
	c.AbortWithStatusJSON(status, gin.H{
		"status": status,
		"error":  err.Error(),
	})
}

// ShouldReturn is a shorthand method for handling resource errors
func ShouldReturn(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	// log error
	logger.Error(err)

	if errors.Is(err, ErrNotFound) || err == store.ErrNotFound {
		NotFound(c, err)
	} else {
		ServerError(c, err)
	}

	return true
}
