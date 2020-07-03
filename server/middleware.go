package server

import (
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/oasishub-indexer/utils/reporting"
	"github.com/gin-gonic/gin"
)

var serverRequestDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
	Namespace: "indexer",
	Subsystem: "oasis.http",
	Name:      "request_duration",
	Desc:      "The total time required to execute http request",
	Tags:      []string{"path"}})

// setupMiddleware sets up middleware for gin application
func (s *Server) setupMiddleware() {
	s.engine.Use(gin.Recovery())
	s.engine.Use(MetricMiddleware())
	s.engine.Use(ErrorReportingMiddleware())
}

// MetricMiddleware is a middleware responsible for logging query execution time metric
func MetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := metrics.NewTimer(serverRequestDuration.WithLabels([]string{c.Request.URL.Path}))
		defer t.ObserveDuration()
		c.Next()
	}
}

func ErrorReportingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer reporting.RecoverError()
		c.Next()
	}
}
