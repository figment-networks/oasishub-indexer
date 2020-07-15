package indexer

import (
	"context"
	"reflect"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	mock "github.com/figment-networks/oasishub-indexer/client/mock"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	timestamppb "github.com/golang/protobuf/ptypes/timestamp"

	"testing"
)

func TestSetup_Run(t *testing.T) {
	setup(t)

	t.Run("returns error when client returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()

		mockClient := mock.NewMockBlockClient(ctrl)
		task := NewHeightMetaRetrieverTask(mockClient)
		pl := &payload{CurrentHeight: 6}

		mockClient.EXPECT().GetByHeight(pl.CurrentHeight).Return(nil, errTestClient).Times(1)

		if result := task.Run(ctx, pl); result != errTestClient {
			t.Errorf("want: %v, got: %v", errTestClient, result)
		}
	})

	tt := struct {
		appVersion   uint64
		blockVersion uint64
		height       int64
		timestamp    *timestamppb.Timestamp
	}{35, 43, 25, ptypes.TimestampNow()}

	t.Run("updates payload.HeightMeta", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx := context.Background()

		mockClient := mock.NewMockBlockClient(ctrl)
		task := NewHeightMetaRetrieverTask(mockClient)

		pl := &payload{CurrentHeight: 6}
		resp := testGetByHeightResp(tt.appVersion, tt.blockVersion, tt.height, tt.timestamp)

		mockClient.EXPECT().GetByHeight(pl.CurrentHeight).Return(resp, nil).Times(1)

		if result := task.Run(ctx, pl); result != nil {
			t.Errorf("want %v; got %v", nil, result)
			return
		}

		expectedHeightMeta := HeightMeta{
			Height:       tt.height,
			Time:         *types.NewTimeFromTimestamp(*tt.timestamp),
			AppVersion:   tt.appVersion,
			BlockVersion: tt.blockVersion,
		}

		if !reflect.DeepEqual(pl.HeightMeta, expectedHeightMeta) {
			t.Errorf("want: %+v, got: %+v", expectedHeightMeta, pl.HeightMeta)
			return
		}
	})

}

func testGetByHeightResp(appVersion, blockVersion uint64, height int64, ts *timestamppb.Timestamp) *blockpb.GetByHeightResponse {
	return &blockpb.GetByHeightResponse{
		Block: &blockpb.Block{
			Header: &blockpb.Header{
				Version: &blockpb.Version{
					App:   appVersion,
					Block: blockVersion,
				},
				Height: height,
				Time:   ts,
			},
		},
	}
}
