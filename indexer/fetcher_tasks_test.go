package indexer

import (
	"context"
	"reflect"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	mock "github.com/figment-networks/oasishub-indexer/client/mock"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"testing"
)

func TestBlockFetcher_Run(t *testing.T) {
	setup(t)

	tests := []struct {
		description   string
		expectedBlock *blockpb.Block
		result        error
	}{
		{"returns error if client errors", nil, testClientErr},
		{"updates payload.RawBlock", testpbBlock(35, 43, 25, ptypes.TimestampNow()), nil},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockClient := mock.NewMockBlockClient(ctrl)
			task := NewBlockFetcherTask(mockClient)

			pl := &payload{CurrentHeight: 20}

			mockClient.EXPECT().GetByHeight(pl.CurrentHeight).Return(&blockpb.GetByHeightResponse{Block: tt.expectedBlock}, tt.result).Times(1)

			if result := task.Run(ctx, pl); result != tt.result {
				t.Errorf("want %v; got %v", tt.result, result)
				return
			}

			// skip payload check if there's an error
			if tt.result != nil {
				return
			}

			if !reflect.DeepEqual(pl.RawBlock, tt.expectedBlock) {
				t.Errorf("want: %+v, got: %+v", tt.expectedBlock, pl.RawBlock)
				return
			}
		})
	}
}

func TestStateFetcher_Run(t *testing.T) {
	setup(t)

	tests := []struct {
		description   string
		expectedState *statepb.State
		result        error
	}{
		{"returns error if client errors", nil, testClientErr},
		{"updates payload.RawState", testpbState(), nil},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockClient := mock.NewMockStateClient(ctrl)
			task := NewStateFetcherTask(mockClient)

			pl := &payload{CurrentHeight: 30}

			mockClient.EXPECT().GetByHeight(pl.CurrentHeight).Return(&statepb.GetByHeightResponse{State: tt.expectedState}, tt.result).Times(1)

			if result := task.Run(ctx, pl); result != tt.result {
				t.Errorf("want %v; got %v", tt.result, result)
				return
			}

			// skip payload check if there's an error
			if tt.result != nil {
				return
			}

			if !reflect.DeepEqual(pl.RawState, tt.expectedState) {
				t.Errorf("want: %+v, got: %+v", tt.expectedState, pl.RawState)
				return
			}
		})
	}
}

func testpbState() *statepb.State {
	return &statepb.State{
		ChainID: randString(10),
		Height:  89,
		Staking: &statepb.Staking{
			TotalSupply: randBytes(10),
			CommonPool:  randBytes(10),
		},
	}
}

func testpbBlock(appVersion, blockVersion uint64, height int64, ts *timestamppb.Timestamp) *blockpb.Block {
	return &blockpb.Block{
		Header: &blockpb.Header{
			Version: &blockpb.Version{
				App:   appVersion,
				Block: blockVersion,
			},
			Height: height,
			Time:   ts,
		},
	}
}
