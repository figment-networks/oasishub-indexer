package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/logger"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

func NewMainSyncerTask(db *store.Store) pipeline.Task {
	return &mainSyncerTask{
		db: db,
	}
}

type mainSyncerTask struct {
	db *store.Store
}

func (t *mainSyncerTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer logTaskDuration(time.Now(), reflect.TypeOf(*t).Name())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("started height=%d", payload.CurrentHeight))

	report, ok := ctx.Value(CtxReport).(*model.Report)
	if !ok {
		return errors.New("report missing in context")
	}

	syncable := &model.Syncable{
		Height: payload.CurrentHeight,
		ReportID: report.ID,
		Time: payload.HeightMeta.Time,
		AppVersion: payload.HeightMeta.AppVersion,
		BlockVersion: payload.HeightMeta.BlockVersion,
		Status: model.SyncableStatusRunning,
		StartedAt: *types.NewTimeFromTime(time.Now()),
	}
	err := t.db.Syncables.CreateOrUpdate(syncable)
	if err != nil {
		return err
	}
	payload.Syncable = syncable
	return nil
}
