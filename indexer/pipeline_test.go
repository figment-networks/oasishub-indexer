package indexer

import (
	mockIndexer "github.com/figment-networks/oasishub-indexer/mock/indexer"
	mockStore "github.com/figment-networks/oasishub-indexer/mock/store"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestIndexingPipeline_getMissingVersionIds(t *testing.T) {
	t.Run("when current index version is higher than smallest for syncable", func(t *testing.T) {
		ctrl:= gomock.NewController(t)
		defer ctrl.Finish()

		syncablesStoreMock := mockStore.NewMockSyncablesStore(ctrl)
		targetsReaderMock := mockIndexer.NewMockTargetsReader(ctrl)

		var smallestIndexVersion int64 = 1
		var currentIndexVersion int64 = 2
		syncablesStoreMock.EXPECT().FindSmallestIndexVersion().Return(&smallestIndexVersion, nil).Times(1)
		targetsReaderMock.EXPECT().GetCurrentVersionID().Return(currentIndexVersion).Times(1)

		s := &store.Store{
			Syncables: syncablesStoreMock,
		}

		testPipeline := indexingPipeline{
			db: s,
			targetsReader: targetsReaderMock,
		}

		ids, err := testPipeline.getMissingVersionIds()
		if err != nil {
			t.Errorf("unexpected error occured; err=%+v", err)
			return
		}

		if len(ids) != 1 {
			t.Errorf("undexpected ids count, want: %d; got: %d", 1, len(ids))
			return
		}

		if ids[0] != 2 {
			t.Errorf("id returned, want: %d; got: %d", ids[0], 2)
			return
		}
	})

	t.Run("when current index version is the same as smallest for syncable", func(t *testing.T) {
		ctrl:= gomock.NewController(t)
		defer ctrl.Finish()

		syncablesStoreMock := mockStore.NewMockSyncablesStore(ctrl)
		targetsReaderMock := mockIndexer.NewMockTargetsReader(ctrl)

		var smallestIndexVersion int64 = 2
		var currentIndexVersion int64 = 2
		syncablesStoreMock.EXPECT().FindSmallestIndexVersion().Return(&smallestIndexVersion, nil).Times(1)
		targetsReaderMock.EXPECT().GetCurrentVersionID().Return(currentIndexVersion).Times(1)

		s := &store.Store{
			Syncables: syncablesStoreMock,
		}

		testPipeline := indexingPipeline{
			db: s,
			targetsReader: targetsReaderMock,
		}

		ids, err := testPipeline.getMissingVersionIds()
		if err != nil {
			t.Errorf("unexpected error occured; err=%+v", err)
			return
		}

		if len(ids) != 2 {
			t.Errorf("getMissingVersionIds undexpected ids count, want: %d; got: %d", 0, len(ids))
		}

		if ids[0] != 1 {
			t.Errorf("id returned, want: %d; got: %d", ids[0], 1)
			return
		}

		if ids[1] != 2 {
			t.Errorf("id returned, want: %d; got: %d", ids[1], 2)
			return
		}
	})

	t.Run("return err when current index version is smaller than smallest for syncable", func(t *testing.T) {
		ctrl:= gomock.NewController(t)
		defer ctrl.Finish()

		syncablesStoreMock := mockStore.NewMockSyncablesStore(ctrl)
		targetsReaderMock := mockIndexer.NewMockTargetsReader(ctrl)

		var smallestIndexVersion int64 = 4
		var currentIndexVersion int64 = 2
		syncablesStoreMock.EXPECT().FindSmallestIndexVersion().Return(&smallestIndexVersion, nil).Times(1)
		targetsReaderMock.EXPECT().GetCurrentVersionID().Return(currentIndexVersion).Times(1)

		s := &store.Store{
			Syncables: syncablesStoreMock,
		}

		testPipeline := indexingPipeline{
			db: s,
			targetsReader: targetsReaderMock,
		}

		_, err := testPipeline.getMissingVersionIds()
		if err == nil {
			t.Errorf("error should be returned")
		}
	})

	t.Run("starts from one when FindSmallestIndexVersion return not found", func(t *testing.T) {
		ctrl:= gomock.NewController(t)
		defer ctrl.Finish()

		syncablesStoreMock := mockStore.NewMockSyncablesStore(ctrl)
		targetsReaderMock := mockIndexer.NewMockTargetsReader(ctrl)

		var currentIndexVersion int64 = 2
		syncablesStoreMock.EXPECT().FindSmallestIndexVersion().Return(nil, store.ErrNotFound).Times(1)
		targetsReaderMock.EXPECT().GetCurrentVersionID().Return(currentIndexVersion).Times(1)

		s := &store.Store{
			Syncables: syncablesStoreMock,
		}

		testPipeline := indexingPipeline{
			db: s,
			targetsReader: targetsReaderMock,
		}

		ids, err := testPipeline.getMissingVersionIds()
		if err != nil {
			t.Errorf("unexpected error occured; err=%+v", err)
			return
		}

		if len(ids) != 2 {
			t.Errorf("undexpected ids count, want: %d; got: %d", 2, len(ids))
			return
		}

		if ids[0] != 1 {
			t.Errorf("id returned, want: %d; got: %d", ids[0], 1)
			return
		}

		if ids[1] != 2 {
			t.Errorf("id returned, want: %d; got: %d", ids[1], 2)
			return
		}
	})
}
