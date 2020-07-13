package indexer

import (
	"context"
	"reflect"
	"time"

	mock "github.com/figment-networks/oasishub-indexer/indexer/mock"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"testing"
)

func TestRun(t *testing.T) {
	setup(t)

	const testReportID types.ID = 64

	t.Run("returns error when report is missing from context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()

		dbMock := mock.NewMockSyncableStore(ctrl)
		task := NewMainSyncerTask(dbMock)

		if result := task.Run(ctx, testSyncerPayload()); result != ErrMissingReportInCtx {
			t.Errorf("want: %v, got: %v", ErrMissingReportInCtx, result)
		}
	})

	t.Run("returns error when db errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := ctxWithReport(testReportID)
		dbErr := errors.New("dberr")

		dbMock := mock.NewMockSyncableStore(ctrl)
		dbMock.EXPECT().CreateOrUpdate(gomock.Any()).Return(dbErr).Times(1)

		task := NewMainSyncerTask(dbMock)

		if result := task.Run(ctx, testSyncerPayload()); result != dbErr {
			t.Errorf("want: %v, got: %v", ErrMissingReportInCtx, result)
		}
	})

	t.Run("updates syncable on payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := ctxWithReport(testReportID)
		payload := testSyncerPayload()
		mockNow := time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
		Now = func() time.Time { return mockNow }

		syncable := &model.Syncable{
			Height:       payload.CurrentHeight,
			ReportID:     testReportID,
			Time:         payload.HeightMeta.Time,
			AppVersion:   payload.HeightMeta.AppVersion,
			BlockVersion: payload.HeightMeta.BlockVersion,
			Status:       model.SyncableStatusRunning,
			StartedAt:    *types.NewTimeFromTime(mockNow),
		}

		dbMock := mock.NewMockSyncableStore(ctrl)
		dbMock.EXPECT().CreateOrUpdate(syncable).Times(1)
		task := NewMainSyncerTask(dbMock)

		if result := task.Run(ctx, payload); result != nil {
			t.Errorf("want %v; got %v", nil, result)
			return
		}

		if !reflect.DeepEqual(payload.Syncable, syncable) {
			t.Errorf("want: %+v, got: %+v", syncable, payload.Syncable)
			return
		}

	})
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
