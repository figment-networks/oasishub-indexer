package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/metric"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
)

const (
	MainSyncerTaskName = "MainSyncer"
)

var (
	ErrMissingReportInCtx = errors.New("report missing in context")
)

func NewMainSyncerTask(db SyncableStore) pipeline.Task {
	return &mainSyncerTask{
		db: db,
	}
}

type mainSyncerTask struct {
	db SyncableStore
}

func (t *mainSyncerTask) GetName() string {
	return MainSyncerTaskName
}

func (t *mainSyncerTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", pipeline.StageSyncer, t.GetName(), payload.CurrentHeight))

	report, ok := ctx.Value(CtxReport).(*model.Report)
	if !ok {
		return ErrMissingReportInCtx
	}

	syncable := &model.Syncable{
		Height:       payload.CurrentHeight,
		ReportID:     report.ID,
		Time:         payload.HeightMeta.Time,
		AppVersion:   payload.HeightMeta.AppVersion,
		BlockVersion: payload.HeightMeta.BlockVersion,
		Status:       model.SyncableStatusRunning,
		StartedAt:    *types.NewTimeFromTime(Now()),
	}
	err := t.db.CreateOrUpdate(syncable)
	if err != nil {
		return err
	}
	payload.Syncable = syncable
	return nil
}
