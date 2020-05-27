package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

var (
	_ pipeline.Sink = (*sink)(nil)
)

func NewSink(db *store.Store) *sink {
	return &sink{
		db: db,
	}
}

type sink struct {
	db *store.Store

	successCount int64
}

func (s *sink) Consume(ctx context.Context, p pipeline.Payload) error {
	payload := p.(*payload)

	logger.DebugJSON(payload,
		logger.Field("process", "pipeline"),
		logger.Field("stage", "sink"),
		logger.Field("height", payload.CurrentHeight),
	)

	payload.Syncable.MarkProcessed()
	if err := s.db.Syncables.Save(payload.Syncable); err != nil {
		return errors.Wrap(err, "failed saving syncable in sink")
	}

	s.successCount += 1

	metric.NumIndexingSuccess.Inc()
	metric.IndexingDuration.Set(payload.Syncable.Duration.Seconds())

	logger.Info(fmt.Sprintf("processing of height %d completed successfully", payload.CurrentHeight))

	return nil
}