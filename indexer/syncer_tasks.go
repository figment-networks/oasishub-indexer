package indexing

import (
	"context"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/pkg/errors"
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
	payload := p.(*payload)

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
