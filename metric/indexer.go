package metric

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	IndexerHeightSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "figment",
		Subsystem: "indexer",
		Name: "height_success",
		Help: "The total number of successfully indexed heights",
	})

	IndexerTotalErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "figment",
		Subsystem: "indexer",
		Name: "total_error",
		Help: "The total number of failures during indexing",
	})

	IndexerHeightDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "figment",
		Subsystem: "indexer",
		Name: "height_duration",
		Help: "The total time required to index one height",
	})

	IndexerTaskDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "figment",
			Subsystem: "indexer",
			Name: "height_task_duration",
			Help: "The total time required to process indexing task",
		},
		[]string{"task"},
	)

	IndexerDbSizeAfterHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "figment",
		Subsystem: "indexer",
		Name: "db_size",
		Help: "The size of the database after indexing of height",
	})
)

// IndexerMetric handles HTTP requests
type IndexerMetric struct {}

// NewIndexerMetric returns a new server instance
func NewIndexerMetric() *IndexerMetric {
	app := &IndexerMetric{}
	return app.init()
}

func (m *IndexerMetric) StartServer(listenAdd string, url string) error {
	logger.Info(fmt.Sprintf("starting metric server at %s...", url), logger.Field("app", "indexer"))

	http.Handle(url, promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	return http.ListenAndServe(listenAdd, nil)
}

func (m *IndexerMetric) init() *IndexerMetric {
	logger.Info("initializing metric server...", logger.Field("app", "indexer"))

	prometheus.MustRegister(IndexerHeightSuccess)
	prometheus.MustRegister(IndexerTotalErrors)
	prometheus.MustRegister(IndexerHeightDuration)
	prometheus.MustRegister(IndexerTaskDuration)
	prometheus.MustRegister(IndexerDbSizeAfterHeight)

	// Add Go module build info.
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	return m
}