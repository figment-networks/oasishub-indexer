package indexer

import "github.com/figment-networks/indexing-engine/metrics"

var (
	indexerTaskDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "oasishub_task",
		Name:      "task_duration",
		Desc:      "The total time required to process indexing task",
		Tags:      []string{"task"},
	})

	indexerTotalErrors = metrics.MustNewCounterWithTags(metrics.Options{
		Namespace: "indexers",
		Subsystem: "oasishub_task",
		Name:      "total_error",
		Desc:      "The total number of failures during indexing",
	})

	indexerHeightSuccess = metrics.MustNewCounterWithTags(metrics.Options{
		Namespace: "indexers",
		Subsystem: "oasishub_task",
		Name:      "height_success",
		Desc:      "The total number of successfully indexed heights",
	}).WithLabels()

	indexerDbSizeAfterHeight = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "oasishub_task",
		Name:      "db_size",
		Desc:      "The size of the database after indexing of height",
	}).WithLabels()

	indexerHeightDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "oasishub_task",
		Name:      "height_duration",
		Desc:      "The total time required to index one height",
	}).WithLabels()
)
