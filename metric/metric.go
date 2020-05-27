package metric

import (
	"fmt"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	NumIndexingSuccess prometheus.Counter
	NumIndexingErr prometheus.Counter
	IndexingDuration prometheus.Gauge
	IndexingTaskDuration *prometheus.GaugeVec
)

// Metric handles HTTP requests
type Metric struct {}

// New returns a new server instance
func New() *Metric {
	app := &Metric{}
	return app.init()
}

func (a *Metric) StartServer(listenAdd string, url string) error {
	logger.Info(fmt.Sprintf("starting server at %s...", url), logger.Field("app", "metric-server"))

	http.Handle(url, promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	return http.ListenAndServe(listenAdd, nil)
}

func (m *Metric) init() *Metric {
	logger.Info("initializing metric server...")

	NumIndexingSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "figment",
		Subsystem: "indexer",
		Name: "height_success",
		Help: "The total number of successfully indexed heights",
	})

	NumIndexingErr = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "figment",
		Subsystem: "indexer",
		Name: "total_error",
		Help: "The total number of failures during indexing",
	})

	IndexingDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "figment",
		Subsystem: "indexer",
		Name: "height_duration",
		Help: "The total time required to index one height",
	})

	IndexingTaskDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "figment",
			Subsystem: "indexer",
			Name: "height_task_duration",
			Help: "The total time required to process indexing task",
		},
		[]string{"task"},
	)

	prometheus.MustRegister(NumIndexingSuccess)
	prometheus.MustRegister(NumIndexingErr)
	prometheus.MustRegister(IndexingDuration)
	prometheus.MustRegister(IndexingTaskDuration)

	// Add Go module build info.
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	return m
}