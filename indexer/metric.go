package indexer

import "github.com/figment-networks/indexing-engine/metrics"

var (
	indexerTaskDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "oasishub.task",
		Name:      "height_task_duration",
		Desc:      "The total time required to process indexing task",
		Tags:      []string{"task"},
	})
)
