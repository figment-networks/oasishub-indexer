package indexer

import "github.com/figment-networks/indexing-engine/metrics"

var (
	indexerDbSizeAfterHeight = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexers",
		Subsystem: "oasishub_task",
		Name:      "db_size",
		Desc:      "The size of the database after indexing of height",
	}).WithLabels()
)
