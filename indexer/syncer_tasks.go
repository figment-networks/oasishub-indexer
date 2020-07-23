package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
)

const (
	TaskNameMainSyncer = "MainSyncer"
)

type SyncerTaskStore interface {
	FindByHeight(height int64) (*model.Syncable, error)
}

func NewMainSyncerTask(db SyncerTaskStore) pipeline.Task {
	return &mainSyncerTask{
		db:             db,
		metricObserver: indexerTaskDuration.WithLabels("TaskNameMainSyncer"),
	}
}

type mainSyncerTask struct {
	db             SyncerTaskStore
	metricObserver metrics.Observer
}

func (t *mainSyncerTask) GetName() string {
	return TaskNameMainSyncer
}

func (t *mainSyncerTask) Run(ctx context.Context, p pipeline.Payload) error {
	timer := metrics.NewTimer(t.metricObserver)
	defer timer.ObserveDuration()

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSyncer, t.GetName(), payload.CurrentHeight))

	syncable, err := t.db.FindByHeight(payload.CurrentHeight)
	if err != nil {
		if err == store.ErrNotFound {
			syncable = &model.Syncable{
				Height:       payload.CurrentHeight,
				Time:         payload.HeightMeta.Time,
				AppVersion:   payload.HeightMeta.AppVersion,
				BlockVersion: payload.HeightMeta.BlockVersion,
				Status:       model.SyncableStatusRunning,
			}
		} else {
			return err
		}
	}

	syncable.StartedAt = *types.NewTimeFromTime(Now())

	report, ok := ctx.Value(CtxReport).(*model.Report)
	if ok {
		syncable.ReportID = report.ID
	}

	payload.Syncable = syncable
	return nil
}
