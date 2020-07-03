package indexing

import (
	"github.com/figment-networks/indexing-engine/metrics"
)

var (
	indexerUseCaseDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "oasishub.task",
		Name:      "usecase_duration",
		Desc:      "The total time required to execute use case",
		Tags:      []string{"task"},
	})
)
