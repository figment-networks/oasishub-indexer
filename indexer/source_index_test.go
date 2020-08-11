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

func TestSource_NewIndexSource(t *testing.T) {
	const versionNum int64 = 1
	const batchSize int64 = 10
	const configStartH int64 = 3

	const zerostartH int64 = 0
	const endH int64 = zerostartH + 5

	setup()

	tests := []struct {
		description string
		startHeight int64

		dbResp *model.Syncable
		dbErr  error

		clientResp *chainpb.GetHeadResponse
		clientErr  error

		expectStartHeight int64
		expectEndHeight   int64
		expectErr         error
	}{
		{description: "should start from last unprocessed block",
			startHeight:       zerostartH,
			dbResp:            testSyncable(zerostartH, false),
			dbErr:             nil,
			clientResp:        testpbChainResp(endH),
			clientErr:         nil,
			expectStartHeight: zerostartH,
			expectEndHeight:   endH,
			expectErr:         nil},
		{description: "should start from next block if last block is already processed",
			startHeight:       zerostartH,
			dbResp:            testSyncable(zerostartH, true),
			dbErr:             nil,
			clientResp:        testpbChainResp(endH),
			clientErr:         nil,
			expectStartHeight: (zerostartH + 1),
			expectEndHeight:   endH,
			expectErr:         nil},
		{description: "handle client errors",
			startHeight:       zerostartH,
			dbResp:            testSyncable(zerostartH, false),
			dbErr:             nil,
			clientResp:        nil,
			clientErr:         errTestClient,
			expectStartHeight: 0,
			expectEndHeight:   0,
			expectErr:         errTestClient},
		{description: "should start from startheight param if > 0",
			startHeight:       3,
			dbResp:            nil,
			dbErr:             store.ErrNotFound,
			clientResp:        testpbChainResp(endH),
			clientErr:         nil,
			expectStartHeight: 3,
			expectEndHeight:   endH,
			expectErr:         nil},
		{description: "should start from config startheight if last block doesnt exist in store",
			startHeight:       zerostartH,
			dbResp:            nil,
			dbErr:             store.ErrNotFound,
			clientResp:        testpbChainResp(endH),
			clientErr:         nil,
			expectStartHeight: configStartH,
			expectEndHeight:   endH,
			expectErr:         nil},
		{description: "handle unexpected db error",
			startHeight:       zerostartH,
			dbResp:            nil,
			dbErr:             errTestDbFind,
			clientResp:        testpbChainResp(endH),
			clientErr:         nil,
			expectStartHeight: 0,
			expectEndHeight:   0,
			expectErr:         errTestDbFind},
		{description: "len should not exceed batch size",
			startHeight:       zerostartH,
			dbResp:            testSyncable(zerostartH, false),
			dbErr:             nil,
			clientResp:        testpbChainResp(batchSize + endH),
			clientErr:         nil,
			expectStartHeight: zerostartH,
			expectEndHeight:   (batchSize - 1),
			expectErr:         nil},
		{description: "error when nothing to process",
			startHeight:       zerostartH,
			dbResp:            testSyncable(endH, false),
			dbErr:             nil,
			clientResp:        testpbChainResp(endH),
			clientErr:         nil,
			expectStartHeight: endH,
			expectEndHeight:   endH,
			expectErr:         ErrNothingToProcess},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			clientMock := mock_client.NewMockChainClient(ctrl)
			clientMock.EXPECT().GetHead().Return(tt.clientResp, tt.clientErr)

			dbMock := mock.NewMockSourceIndexStore(ctrl)
			dbMock.EXPECT().FindMostRecent().Return(tt.dbResp, tt.dbErr)

			cfg := &config.Config{FirstBlockHeight: configStartH}
			source, err := NewIndexSource(cfg, dbMock, clientMock, tt.startHeight, batchSize)

			if err != tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
			}

			if tt.expectErr != nil {
				// skip rest of tests if error expected
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

func testpbChainResp(height int64) *chainpb.GetHeadResponse {
	return &chainpb.GetHeadResponse{
		Height: height,
	}
}
