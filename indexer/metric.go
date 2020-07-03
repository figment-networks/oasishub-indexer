package indexer

import "github.com/figment-networks/indexing-engine/metrics"

var (
	indexerTaskDuration = metrics.MustNewGaugeWithTags(metrics.Options{
		Namespace: "figment",
		Subsystem: "indexer",
		Name:      "height_task_duration",
		Desc:      "The total time required to process indexing task",
		Tags:      []string{"task"},
	})
)
