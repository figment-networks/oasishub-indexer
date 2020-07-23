package indexer

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/figment-networks/oasishub-indexer/config"
	mock "github.com/figment-networks/oasishub-indexer/mock/indexer"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/golang/mock/gomock"
)

func TestSource_NewBackfillSource(t *testing.T) {
	const indexVersion int64 = 6
	const startH int64 = 0
	const endH int64 = startH + 5

	setup(t)

	tests := []struct {
		description string
		dbErr       error
		expectErr   error
	}{
		{description: "should return nil if no database errors",
			dbErr:     nil,
			expectErr: nil,
		},
		{description: "should return error if unexpected database error",
			dbErr:     errTestDbFind,
			expectErr: errTestDbFind,
		},
		{description: "should return error if no syncable in store",
			dbErr:     store.ErrNotFound,
			expectErr: ErrNothingToBackfill},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v when fetching start", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)

			dbMock := mock.NewMockBackfillSourceStore(ctrl)
			if tt.dbErr != nil {
				dbMock.EXPECT().FindFirstByDifferentIndexVersion(indexVersion).Return(nil, tt.dbErr)
			} else {
				dbMock.EXPECT().FindFirstByDifferentIndexVersion(indexVersion).Return(&model.Syncable{
					Height: startH,
				}, nil)
			}

			dbMock.EXPECT().FindMostRecentByDifferentIndexVersion(indexVersion).Return(&model.Syncable{
				Height: endH,
			}, nil)

			cfg := &config.Config{}
			srcCgf := &BackfillSourceConfig{indexVersion}
			source, err := NewBackfillSource(cfg, dbMock, srcCgf)

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
			}

			if tt.expectErr != nil {
				// skip rest of tests if error expected
				return
			}

			expectLen := (endH - startH + 1)

			if expectLen != source.Len() {
				t.Errorf("unexpected source.len, want %v; got %v", expectLen, source.Len())
				return
			}

			ctx := context.Background()
			pl := &payload{}
			expectCurrent := startH

			for i := 1; i < int(expectLen); i++ {
				if expectCurrent != source.Current() {
					t.Errorf("unexpected source.current, want %v; got %v", expectCurrent, source.Current())
				}

				if ok := source.Next(ctx, pl); !ok {
					t.Errorf("expected source.Next to return true on call %v", i)
				}
				expectCurrent++ // current height should increase by 1 after each run
			}

			//should be no more runs left
			if ok := source.Next(ctx, pl); ok {
				t.Errorf("expected source.Next to return false on call %v", expectLen)
			}
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v when fetching most recent", tt.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)

			dbMock := mock.NewMockBackfillSourceStore(ctrl)
			dbMock.EXPECT().FindFirstByDifferentIndexVersion(indexVersion).Return(&model.Syncable{
				Height: startH,
			}, nil)

			if tt.dbErr != nil {
				dbMock.EXPECT().FindMostRecentByDifferentIndexVersion(indexVersion).Return(nil, tt.dbErr)
			} else {
				dbMock.EXPECT().FindMostRecentByDifferentIndexVersion(indexVersion).Return(&model.Syncable{
					Height: endH,
				}, nil)
			}

			cfg := &config.Config{}
			srcCgf := &BackfillSourceConfig{indexVersion}
			_, err := NewBackfillSource(cfg, dbMock, srcCgf)

			if !errors.Is(err, tt.expectErr) {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
			}
		})
	}
}
