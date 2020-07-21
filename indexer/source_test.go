package indexer

import (
	"context"
	"testing"
	"time"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/chain/chainpb"
	"github.com/figment-networks/oasishub-indexer/config"
	mock_client "github.com/figment-networks/oasishub-indexer/mock/client"
	mock "github.com/figment-networks/oasishub-indexer/mock/indexer"
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
		expectEndHeight   int64
		expectErr         error
	}{
		{"should start from last unprocessed block", testSyncable(startH, false), nil, testpbChainResp(endH), nil, startH, endH, nil},
		{"should start from next block if last block is already processed", testSyncable(startH, true), nil, testpbChainResp(endH), nil, (startH + 1), endH, nil},
		{"handle client errors", testSyncable(startH, false), nil, nil, errTestClient, 0, 0, errTestClient},
		{"should start from config startheight if last block doesnt exist in store", nil, store.ErrNotFound, testpbChainResp(endH), nil, configStartH, endH, nil},
		{"handle unexpected db error", nil, errTestDb, testpbChainResp(endH), nil, 0, 0, errTestDb},
		{"len should not exceed batch size", testSyncable(startH, false), nil, testpbChainResp(batchSize + endH), nil, startH, (batchSize - 1), nil},
		{"error when nothing to process", testSyncable(endH, false), nil, testpbChainResp(endH), nil, endH, endH, ErrNothingToProcess},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			clientMock := mock_client.NewMockChainClient(ctrl)
			clientMock.EXPECT().GetHead().Return(tt.clientResp, tt.clientErr)

			dbMock := mock.NewMockSyncableStore(ctrl)
			dbMock.EXPECT().FindMostRecent().Return(tt.dbResp, tt.dbErr)

			cfg := &config.Config{FirstBlockHeight: configStartH}
			source := NewSource(cfg, dbMock, clientMock, versionNum, batchSize)

			if tt.expectErr != source.Err() {
				t.Errorf("unexpected source.err, want %v; got %v", tt.expectErr, source.Err())
				return
			}

			expectLen := (tt.expectEndHeight - tt.expectStartHeight + 1)

			if expectLen != source.Len() {
				t.Errorf("unexpected source.len, want %v; got %v", expectLen, source.Len())
				return
			}

			ctx := context.Background()
			pl := &payload{}
			expectCurrent := tt.expectStartHeight

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
