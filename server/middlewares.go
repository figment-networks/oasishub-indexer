package server

import (
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/utils/reporting"
	"github.com/gin-gonic/gin"
	"time"
)

// setupMiddlewares sets up middleware for gin application
func (s *Server) setupMiddlewares() {
	s.engine.Use(gin.Recovery())
	s.engine.Use(MetricMiddleware())
	s.engine.Use(ErrorReportingMiddleware())
}

// MetricMiddleware is a middleware responsible for logging query execution time metric
func MetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		elapsed := time.Since(t)

		metric.DatabaseQueryDuration.WithLabelValues(c.Request.URL.Path).Set(elapsed.Seconds())
	}
}

func ErrorReportingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer reporting.RecoverError()
		c.Next()
	}
}
