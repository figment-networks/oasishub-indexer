package indexing

import (
	"github.com/figment-networks/indexing-engine/metrics"
)

var (
	indexerUseCaseDuration = metrics.MustNewGaugeWithTags(metrics.Options{
		Namespace: "figment",
		Subsystem: "indexer",
		Name:      "use_case_duration",
		Desc:      "The total time required to execute use case",
		Tags:      []string{"task"},
	})
)
