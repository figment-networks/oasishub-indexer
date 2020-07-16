package indexer

import "github.com/figment-networks/indexing-engine/metrics"

var (
	indexerTaskDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "oasishub_task",
		Name:      "height_task_duration",
		Desc:      "The total time required to process indexing task",
		Tags:      []string{"task"},
	})

	indexerTotalErrors = metrics.MustNewCounterWithTags(metrics.Options{
		Namespace: "indexers",
		Subsystem: "oasishub_task",
		Name:      "total_error",
		Desc:      "The total number of failures during indexing",
	})
)
