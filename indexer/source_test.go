package indexer

import (
	"context"
	"testing"
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	climock "github.com/figment-networks/oasishub-indexer/client/mock"
	"github.com/figment-networks/oasishub-indexer/config"
	mock "github.com/figment-networks/oasishub-indexer/indexer/mock"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
)

func TestSource_NewSource(t *testing.T) {
	const versionNum int64 = 1
	const batchSize int64 = 10
	const configStartH int64 = 3
	const startH int64 = 0
	const endH int64 = startH + 5

	setup(t)

	tests := []struct {
		description string

		dbResp *model.Syncable
		dbErr  error

		clientResp *chainpb.GetCurrentResponse
		clientErr  error

		expectStartHeight int64
		expectLen         int64
	}{
		{"should start from last unprocessed block", testSyncable(startH, false), nil, testpbChainResp(endH), nil, startH, (endH - startH + 1)},
		{"should start from next block if last block is already processed", testSyncable(startH, true), nil, testpbChainResp(endH), nil, (startH + 1), (endH - (startH + 1) + 1)},
		{"handle client errors", testSyncable(startH, false), nil, nil, errTestClient, 0, 1},
		{"should start from config startheight if last block doesnt exist in store", nil, store.ErrNotFound, testpbChainResp(endH), nil, configStartH, (endH - configStartH + 1)},
		{"handle unexpected db error", nil, errTestDb, testpbChainResp(endH), nil, 0, 1},
		{"len should not exceed batch size", testSyncable(startH, false), nil, testpbChainResp(batchSize + endH), nil, startH, (batchSize)},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			clientMock := climock.NewMockChainClient(ctrl)
			clientMock.EXPECT().GetHead().Return(tt.clientResp, tt.clientErr)

			dbMock := mock.NewMocksourceStore(ctrl)
			dbMock.EXPECT().FindMostRecent().Return(tt.dbResp, tt.dbErr)

			cfg := &config.Config{FirstBlockHeight: configStartH}
			source := NewSource(cfg, dbMock, clientMock, versionNum, batchSize)

			if tt.clientErr != nil && source.Err() != tt.clientErr {
				t.Errorf("unexpected source.err, want %v; got %v", tt.clientErr, source.Err())
			}

			if tt.dbErr != nil && tt.dbErr != store.ErrNotFound && source.Err() != tt.dbErr {
				t.Errorf("unexpected source.err, want %v; got %v", tt.dbErr, source.Err())
			}

			if tt.expectLen != source.Len() {
				t.Errorf("unexpected source.len, want %v; got %v", tt.expectLen, source.Len())
			}

			ctx := context.Background()
			pl := &payload{}
			expectCurrent := tt.expectStartHeight

			for i := 0; i < int(tt.expectLen-1); i++ {
				if expectCurrent != source.Current() {
					t.Errorf("unexpected source.current, want %v; got %v", expectCurrent, source.Current())
				}

				if ok := source.Next(ctx, pl); !ok {
					t.Errorf("unexpected number of runs, want %v; got %v", tt.expectLen, i)
				}
				expectCurrent++ // current height should increase by 1 after each run
			}

			//should be no more runs left
			if ok := source.Next(ctx, pl); ok {
				t.Errorf("unexpected number of runs, want %v; got %v", tt.expectLen, (tt.expectLen + 1))
			}
		})
	}
}

func testSyncable(height int64, processed bool) *model.Syncable {
	c := &model.Syncable{
		Height: height,
	}
	if processed {
		c.ProcessedAt = types.NewTimeFromTime(time.Now())
	}
	return c
}

func testpbChainResp(height int64) *chainpb.GetCurrentResponse {
	return &chainpb.GetCurrentResponse{
		Chain: testpbChain(height),
	}
}

func testpbChain(height int64) *chainpb.Chain {
	return &chainpb.Chain{
		Height: height,
	}
}
