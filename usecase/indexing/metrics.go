package indexing

import (
	"github.com/figment-networks/indexing-engine/metrics"
)

var (
	indexerUseCaseDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "figment",
		Subsystem: "indexer",
		Name:      "use_case_duration",
		Desc:      "The total time required to execute use case",
		Tags:      []string{"task"},
	})
)
