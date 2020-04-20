package startpipeline

import (
	"context"
	mock_syncablerepo "github.com/figment-networks/oasishub-indexer/mock/repos/syncablerepo"
	"github.com/figment-networks/oasishub-indexer/models/syncable"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/figment-networks/oasishub-indexer/utils/errors"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func Test_Syncer(t *testing.T) {
	startH := types.Height(1)
	endH := types.Height(10)
	pld := &payload{
		StartHeight:   startH,
		EndHeight:     endH,
		CurrentHeight: startH,
		RetrievedAt:   *types.NewTimeFromTime(time.Now()),
		BlockSyncable: &syncable.Model{
			Type: syncable.BlockType,
		},
		StateSyncable: &syncable.Model{
			Type: syncable.StateType,
		},
		ValidatorsSyncable: &syncable.Model{
			Type: syncable.ValidatorsType,
		},
		TransactionsSyncable: &syncable.Model{
			Type: syncable.TransactionsType,
		},
	}
	ctx := context.Background()

	t.Run("Consume() works as expected when syncables are not already in database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := syncable.Model{}
		syncableDbRepo := mock_syncablerepo.NewMockDbRepo(ctrl)
		syncableDbRepo.EXPECT().
			Create(gomock.Any()).
			Return(nil).
			Times(4)
		syncableDbRepo.EXPECT().
			GetByHeight(gomock.Any(), gomock.Eq(startH)).
			Return(nil, errors.NewErrorFromMessage("not found", errors.NotFoundError)).
			Times(4)

		syncableProxyRepo := mock_syncablerepo.NewMockProxyRepo(ctrl)
		syncableProxyRepo.EXPECT().
			GetByHeight(gomock.Any(), gomock.Eq(startH)).
			Return(&s, nil).
			Times(4)

		syncer := NewSyncer(syncableDbRepo, syncableProxyRepo)

		updatedPayload, err := syncer.Process(ctx, pld)
		p := updatedPayload.(*payload)

		if err != nil {
			t.Errorf("should not return error. Err: %v", err)
		}

		if p.BlockSyncable == nil {
			t.Errorf("payload.BlockSyncable should be set")
		}

		if p.StateSyncable == nil {
			t.Errorf("payload.StateSyncable should be set")
		}

		if p.ValidatorsSyncable == nil {
			t.Errorf("payload.ValidatorsSyncable should be set")
		}

		if p.TransactionsSyncable == nil {
			t.Errorf("payload.TransactionsSyncable should be set")
		}
	})

	t.Run("Consume() works as expected when syncables are already in database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := syncable.Model{}
		syncableDbRepo := mock_syncablerepo.NewMockDbRepo(ctrl)
		syncableDbRepo.EXPECT().
			Create(gomock.Any()).
			Return(nil).
			Times(0)
		syncableDbRepo.EXPECT().
			GetByHeight(gomock.Any(), gomock.Eq(startH)).
			Return(&s, nil).
			Times(4)

		syncableProxyRepo := mock_syncablerepo.NewMockProxyRepo(ctrl)
		syncableProxyRepo.EXPECT().
			GetByHeight(gomock.Any(), gomock.Eq(startH)).
			Return(&s, nil).
			Times(0)

		syncer := NewSyncer(syncableDbRepo, syncableProxyRepo)

		updatedPayload, err := syncer.Process(ctx, pld)
		p := updatedPayload.(*payload)

		if err != nil {
			t.Errorf("should not return error. Err: %v", err)
		}

		if p.BlockSyncable == nil {
			t.Errorf("payload.BlockSyncable should be set")
		}

		if p.StateSyncable == nil {
			t.Errorf("payload.StateSyncable should be set")
		}

		if p.ValidatorsSyncable == nil {
			t.Errorf("payload.ValidatorsSyncable should be set")
		}

		if p.TransactionsSyncable == nil {
			t.Errorf("payload.TransactionsSyncable should be set")
		}
	})
}
