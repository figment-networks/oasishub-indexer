package indexer

import (
	mock_store "github.com/figment-networks/oasishub-indexer/mock/store"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestPipelineStatusChecker_getStatus(t *testing.T) {
	tests := []struct {
		description               string
		currentIndexVersion       int64
		setStore                  func(*mock_store.MockSyncablesStore)
		expectError               bool
		expectedIsUpToDate        bool
		expectedIsPristine		  bool
		expectedMissingVersionIds []int64
	}{
		{
			description:         "when smallest and current index versions are the same, it is up to date and returns versions 1 and 2",
			currentIndexVersion: 2,
			setStore: func(mock *mock_store.MockSyncablesStore) {
				var smallestIndexVersion int64 = 2
				mock.EXPECT().FindSmallestIndexVersion().Return(&smallestIndexVersion, nil).Times(1)
			},
			expectError:               false,
			expectedIsUpToDate:        true,
			expectedIsPristine:        false,
			expectedMissingVersionIds: []int64{1, 2},
		},
		{
			description:         "when smallest is smaller than current index versions, it is not up to date and returns versions 3 and 4",
			currentIndexVersion: 4,
			setStore: func(mock *mock_store.MockSyncablesStore) {
				var smallestIndexVersion int64 = 2
				mock.EXPECT().FindSmallestIndexVersion().Return(&smallestIndexVersion, nil).Times(1)
			},
			expectError:               false,
			expectedIsUpToDate:        false,
			expectedIsPristine:        false,
			expectedMissingVersionIds: []int64{3, 4},
		},
		{
			description:         "when smallest is larger than current index versions, it return error",
			currentIndexVersion: 3,
			setStore: func(mock *mock_store.MockSyncablesStore) {
				var smallestIndexVersion int64 = 4
				mock.EXPECT().FindSmallestIndexVersion().Return(&smallestIndexVersion, nil).Times(1)
			},
			expectError: true,
		},
		{
			description:         "when FindSmallestIndexVersion return not found error, it return all version ids",
			currentIndexVersion: 2,
			setStore: func(mock *mock_store.MockSyncablesStore) {
				mock.EXPECT().FindSmallestIndexVersion().Return(nil, store.ErrNotFound).Times(1)
			},
			expectError:               false,
			expectedIsUpToDate:        false,
			expectedIsPristine:        true,
			expectedMissingVersionIds: []int64{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			syncableStoreMock := mock_store.NewMockSyncablesStore(ctrl)

			tt.setStore(syncableStoreMock)

			checker := pipelineStatusChecker{
				currentIndexVersion: tt.currentIndexVersion,
				store:               syncableStoreMock,
			}

			status, err := checker.getStatus()
			if err != nil {
				if !tt.expectError {
					t.Errorf("getStatus() should not return error, got: %v", err)
				}
				return
			} else {
				if tt.expectError {
					t.Errorf("getStatus() should return error")
					return
				}
			}

			if status.isUpToDate != tt.expectedIsUpToDate {
				t.Errorf("unexpected isUpToDate, want: %v; got: %v", tt.expectedIsUpToDate, status.isUpToDate)
			}

			if status.isPristine != tt.expectedIsPristine {
				t.Errorf("unexpected isPristine, want: %v; got: %v", tt.expectedIsPristine, status.isPristine)
			}

			if len(tt.expectedMissingVersionIds) != len(status.missingVersionIds) {
				t.Errorf("unexpected missingVersionIds size, want: %d; got: %d", len(tt.expectedMissingVersionIds), len(status.missingVersionIds))
			}

			for i, gotId := range status.missingVersionIds {
				wantId := tt.expectedMissingVersionIds[i]
				if wantId != gotId {
					t.Errorf("unexpecte missing version id, want: %d; got: %d", wantId, gotId)
				}
			}
		})

	}
}
