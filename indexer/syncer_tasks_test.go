package indexer

import (
	"context"
	"reflect"
	"time"

	mock "github.com/figment-networks/oasishub-indexer/mock/indexer"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"testing"
)

func TestSyncer_Run(t *testing.T) {
	const testReportID types.ID = 64
	dbErr := errors.New("dberr")

	t.Run("returns error when db errors", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		pl := testSyncerPayload()

		dbMock := mock.NewMockSyncerTaskStore(ctrl)
		dbMock.EXPECT().FindByHeight(pl.CurrentHeight).Return(nil, dbErr).Times(1)

		task := NewMainSyncerTask(dbMock)

		if err := task.Run(ctx, pl); err != dbErr {
			t.Errorf("unexpected error, want: %v, got: %v", dbErr, err)
		}
	})

	t.Run("updates existing syncable", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		ctx := context.Background()
		pl := testSyncerPayload()

		expectSyncable := &model.Syncable{
			Height:       pl.CurrentHeight,
			Time:         pl.HeightMeta.Time,
			AppVersion:   pl.HeightMeta.AppVersion,
			BlockVersion: pl.HeightMeta.BlockVersion,
			Status:       model.SyncableStatusRunning,
			StartedAt:    *types.NewTimeFromTime(time.Now()),
		}

		dbMock := mock.NewMockSyncerTaskStore(ctrl)
		dbMock.EXPECT().FindByHeight(pl.CurrentHeight).Return(expectSyncable, nil).Times(1)
		task := NewMainSyncerTask(dbMock)

		if err := task.Run(ctx, pl); err != nil {
			t.Errorf("unexpected error, want %v; got %v", nil, err)
			return
		}

		if !reflect.DeepEqual(pl.Syncable, expectSyncable) {
			t.Errorf("unexpected payload.Syncable, want: %+v, got: %+v", expectSyncable, pl.Syncable)
			return
		}
	})

	t.Run("updates new syncable", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		ctx := context.Background()
		pl := testSyncerPayload()

		dbMock := mock.NewMockSyncerTaskStore(ctrl)
		dbMock.EXPECT().FindByHeight(pl.CurrentHeight).Return(nil, store.ErrNotFound).Times(1)
		task := NewMainSyncerTask(dbMock)

		if err := task.Run(ctx, pl); err != nil {
			t.Errorf("unexpected error, want %v; got %v", nil, err)
			return
		}

		if pl.Syncable.Height != pl.CurrentHeight ||
			pl.Syncable.Time != pl.HeightMeta.Time ||
			pl.Syncable.AppVersion != pl.HeightMeta.AppVersion ||
			pl.Syncable.BlockVersion != pl.HeightMeta.BlockVersion ||
			pl.Syncable.Status != model.SyncableStatusRunning ||
			pl.Syncable.StartedAt.IsZero() {
			t.Errorf("unexpected payload.Syncable %+v", pl.Syncable)
			return
		}
	})

	t.Run("Adds reportId to syncable if it exists in context", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		ctx := ctxWithReport(testReportID)
		pl := testSyncerPayload()

		dbMock := mock.NewMockSyncerTaskStore(ctrl)
		dbMock.EXPECT().FindByHeight(pl.CurrentHeight).Return(&model.Syncable{}, nil).Times(1)
		task := NewMainSyncerTask(dbMock)

		if err := task.Run(ctx, pl); err != nil {
			t.Errorf("unexpected error, want: %v, got: %v", nil, err)
		}

		if pl.Syncable.ReportID != testReportID {
			t.Errorf("unexpected reportID in syncable, want: %v, got: %v", testReportID, pl.Syncable.ReportID)
		}
	})
}

func ctxWithReport(modelID types.ID) context.Context {
	ctx := context.Background()
	report := &model.Report{
		Model: &model.Model{ID: modelID},
	}

	return context.WithValue(ctx, CtxReport, report)
}

func testSyncerPayload() *payload {
	return &payload{
		CurrentHeight: 10,
		HeightMeta: HeightMeta{
			Time:         *types.NewTimeFromTime(time.Now()),
			AppVersion:   4,
			BlockVersion: 6,
		},
	}
}
