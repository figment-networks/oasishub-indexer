package indexer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/figment-networks/oasishub-indexer/config"
	mock_indexer "github.com/figment-networks/oasishub-indexer/mock/indexer"
	"github.com/figment-networks/oasishub-indexer/model"
	"github.com/figment-networks/oasishub-indexer/store"
	"github.com/figment-networks/oasishub-indexer/types"
	"github.com/golang/mock/gomock"
)

const (
	testHeight           = 17
	testValidatorAddress = "test_address"
)

var (
	ErrValidatorSeqFindByHeight = errors.New("could not find test")
	ErrCouldNotFindByAddress    = errors.New("could not find test")

	testCfg = &config.Config{
		FirstBlockHeight: 1,
	}
)

func TestSystemEventCreatorTask_Run(t *testing.T) {
	t.Run("creates system events", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		validatorSeqStoreMock := mock_indexer.NewMockSystemEventCreatorStore(ctrl)

		MaxValidatorSequences = 10
		MissedInRowThreshold = 5
		MissedForMaxThreshold = 5

		prevHeightValidatorSequences := []model.ValidatorSeq{
			newValidatorSeq(testValidatorAddress, 1000, 0, true),
		}
		lastValidatorSeqsForValidator1 := []model.ValidatorSeq{
			newValidatorSeq(testValidatorAddress, 1000, 0, false),
			newValidatorSeq(testValidatorAddress, 1000, 0, false),
			newValidatorSeq(testValidatorAddress, 1000, 0, false),
			newValidatorSeq(testValidatorAddress, 1000, 0, false),
			newValidatorSeq(testValidatorAddress, 1000, 0, false),
			newValidatorSeq(testValidatorAddress, 1000, 0, false),
			newValidatorSeq(testValidatorAddress, 1000, 0, true),
			newValidatorSeq(testValidatorAddress, 1000, 0, true),
			newValidatorSeq(testValidatorAddress, 1000, 0, true),
			newValidatorSeq(testValidatorAddress, 1000, 0, true),
		}
		lastValidatorSeqsForValidator2 := []model.ValidatorSeq{
			newValidatorSeq("validator_address_1", 1000, 0, true),
			newValidatorSeq("validator_address_1", 1000, 0, true),
			newValidatorSeq("validator_address_1", 1000, 0, true),
			newValidatorSeq("validator_address_1", 1000, 0, true),
		}
		payload := testPayload()
		payload.NewValidatorSequences = []model.ValidatorSeq{
			newValidatorSeq(testValidatorAddress, 100000, 0, false),
		}
		payload.UpdatedValidatorSequences = []model.ValidatorSeq{
			newValidatorSeq("validator_address_1", 1000, 0, false),
		}

		validatorSeqStoreMock.EXPECT().FindByHeight(gomock.Any()).Return(prevHeightValidatorSequences, nil).Times(1)
		gomock.InOrder(
			validatorSeqStoreMock.EXPECT().FindLastByAddress(gomock.Any(), gomock.Any()).Return(lastValidatorSeqsForValidator1, nil),
			validatorSeqStoreMock.EXPECT().FindLastByAddress(gomock.Any(), gomock.Any()).Return(lastValidatorSeqsForValidator2, nil),
		)

		task := NewSystemEventCreatorTask(testCfg, validatorSeqStoreMock)
		err := task.Run(ctx, payload)
		if err != nil {
			t.Errorf("unexpected result, run should not return error: %v", err)
			return
		}

		if len(payload.SystemEvents) == 0 {
			t.Errorf("there should be system events added to the payload")
		}
	})

	t.Run("FindByHeight returns error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		validatorSeqStoreMock := mock_indexer.NewMockSystemEventCreatorStore(ctrl)

		payload := testPayload()

		validatorSeqStoreMock.EXPECT().FindByHeight(gomock.Any()).Return(nil, ErrValidatorSeqFindByHeight).Times(1)

		task := NewSystemEventCreatorTask(testCfg, validatorSeqStoreMock)
		err := task.Run(ctx, payload)
		if err == nil {
			t.Errorf("unexpected result, run should return error")
		}
	})

	t.Run("FindByHeight returns not found error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		validatorSeqStoreMock := mock_indexer.NewMockSystemEventCreatorStore(ctrl)

		lastValidatorSeqsForValidator1 := []model.ValidatorSeq{
			newValidatorSeq(testValidatorAddress, 1000, 0, false),
		}
		payload := testPayload()
		payload.NewValidatorSequences = []model.ValidatorSeq{
			newValidatorSeq(testValidatorAddress, 100000, 0, false),
		}

		validatorSeqStoreMock.EXPECT().FindByHeight(gomock.Any()).Return(nil, store.ErrNotFound).Times(1)
		validatorSeqStoreMock.EXPECT().FindLastByAddress(gomock.Any(), gomock.Any()).Return(lastValidatorSeqsForValidator1, nil).Times(1)

		task := NewSystemEventCreatorTask(testCfg, validatorSeqStoreMock)
		err := task.Run(ctx, payload)
		if err != nil {
			t.Errorf("unexpected result, run should not return error: %v", err)
			return
		}

		if len(payload.SystemEvents) == 0 {
			t.Errorf("there should be system events added to the payload")
		}
	})

	t.Run("FindLastByAddress returns error", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		validatorSeqStoreMock := mock_indexer.NewMockSystemEventCreatorStore(ctrl)

		prevHeightValidatorSequences := []model.ValidatorSeq{
			newValidatorSeq(testValidatorAddress, 1000, 0, true),
		}
		payload := testPayload()
		payload.NewValidatorSequences = []model.ValidatorSeq{
			newValidatorSeq(testValidatorAddress, 100000, 0, false),
		}

		validatorSeqStoreMock.EXPECT().FindByHeight(gomock.Any()).Return(prevHeightValidatorSequences, nil).Times(1)
		validatorSeqStoreMock.EXPECT().FindLastByAddress(gomock.Any(), gomock.Any()).Return(nil, ErrCouldNotFindByAddress).Times(1)

		task := NewSystemEventCreatorTask(testCfg, validatorSeqStoreMock)
		err := task.Run(ctx, payload)
		if err == nil {
			t.Errorf("unexpected result, run should return error")
		}
	})
}

func TestSystemEventCreatorTask_getValueChangeSystemEvents(t *testing.T) {
	tests := []struct {
		description                   string
		activeEscrowBalanceChangeRate float64
		commissionChangeRate          float64
		expectedCount                 int
		expectedKind                  model.SystemEventKind
	}{
		{"returns no system events when active escrow balance haven't changed", 0, 0, 0, ""},
		{"returns no system events when active escrow balance change smaller than 0.1", 0.09, 0, 0, ""},
		{"returns one activeEscrowBalanceChange1 system event when active escrow balance change is 0.1", 0.1, 0, 1, model.SystemEventActiveEscrowBalanceChange1},
		{"returns one activeEscrowBalanceChange1 system events when active escrow balance change is 0.9", 0.9, 0, 1, model.SystemEventActiveEscrowBalanceChange1},
		{"returns one activeEscrowBalanceChange2 system events when active escrow balance change is 1", 1, 0, 1, model.SystemEventActiveEscrowBalanceChange2},
		{"returns one activeEscrowBalanceChange2 system events when active escrow balance change is 9", 9, 0, 1, model.SystemEventActiveEscrowBalanceChange2},
		{"returns one activeEscrowBalanceChange3 system events when active escrow balance change is 10", 10, 0, 1, model.SystemEventActiveEscrowBalanceChange3},
		{"returns one activeEscrowBalanceChange3 system events when active escrow balance change is 100", 100, 0, 1, model.SystemEventActiveEscrowBalanceChange3},
		{"returns one activeEscrowBalanceChange3 system events when active escrow balance change is 200", 200, 0, 1, model.SystemEventActiveEscrowBalanceChange3},

		{"returns no system events when commission haven't changed", 0, 0, 0, ""},
		{"returns no system events when commission change smaller than 0.1", 0, 0.09, 0, ""},
		{"returns one commissionChange1 system event when commission change is 0.1", 0, 0.1, 1, model.SystemEventCommissionChange1},
		{"returns one commissionChange1 system events when commission change is 0.9", 0, 0.9, 1, model.SystemEventCommissionChange1},
		{"returns one commissionChange2 system events when commission change is 1", 0, 1, 1, model.SystemEventCommissionChange2},
		{"returns one commissionChange2 system events when commission change is 9", 0, 9, 1, model.SystemEventCommissionChange2},
		{"returns one commissionChange3 system events when commission change is 10", 0, 10, 1, model.SystemEventCommissionChange3},
		{"returns one commissionChange3 system events when commission change is 100", 0, 100, 1, model.SystemEventCommissionChange3},
		{"returns one commissionChange3 system events when commission change is 200", 0, 200, 1, model.SystemEventCommissionChange3},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validatorSeqStoreMock := mock_indexer.NewMockSystemEventCreatorStore(ctrl)

			var activeEscrowBalanceBefore int64 = 1000
			activeEscrowBalanceAfter := float64(activeEscrowBalanceBefore) + (float64(activeEscrowBalanceBefore) * tt.activeEscrowBalanceChangeRate / 100)

			var commissionBefore int64 = 1000
			commissionAfter := float64(commissionBefore) + (float64(commissionBefore) * tt.commissionChangeRate / 100)

			prevHeightValidatorSequences := []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, activeEscrowBalanceBefore, commissionBefore, true),
			}
			currHeightValidatorSequences := []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, int64(activeEscrowBalanceAfter), int64(commissionAfter), true),
			}

			task := NewSystemEventCreatorTask(testCfg, validatorSeqStoreMock)
			createdSystemEvents, _ := task.getValueChangeSystemEvents(currHeightValidatorSequences, prevHeightValidatorSequences)

			if len(createdSystemEvents) != tt.expectedCount {
				t.Errorf("unexpected system event count, want %v; got %v", tt.expectedCount, len(createdSystemEvents))
				return
			}

			if len(createdSystemEvents) > 0 && createdSystemEvents[0].Kind != tt.expectedKind {
				t.Errorf("unexpected system event kind, want %v; got %v", tt.expectedKind, createdSystemEvents[0].Kind)
			}
		})
	}
}

func TestSystemEventCreatorTask_getActiveSetPresenceChangeSystemEvents(t *testing.T) {
	tests := []struct {
		description    string
		prevHeightList []model.ValidatorSeq
		currHeightList []model.ValidatorSeq
		expectedCount  int
		expectedKinds  []model.SystemEventKind
	}{
		{
			description: "returns no system events when validator is both in prev and current lists",
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			expectedCount: 0,
		},
		{
			description:    "returns no system events when validator is not in prev nor current lists",
			prevHeightList: []model.ValidatorSeq{},
			currHeightList: []model.ValidatorSeq{},
			expectedCount:  0,
		},
		{
			description:    "returns one joined_active_set system events when validator is not in prev and is in current lists",
			prevHeightList: []model.ValidatorSeq{},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventJoinedActiveSet},
		},
		{
			description: "returns one left_active_set system events when validator is in prev and is not in current lists",
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{},
			expectedCount:  1,
			expectedKinds:  []model.SystemEventKind{model.SystemEventLeftActiveSet},
		},
		{
			description: "returns 2 joined_active_set system events when validators are in prev and not in current lists",
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
				newValidatorSeq("address1", 1000, 0, true),
				newValidatorSeq("address2", 1000, 0, true),
			},
			expectedCount: 2,
			expectedKinds: []model.SystemEventKind{model.SystemEventJoinedActiveSet, model.SystemEventJoinedActiveSet},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validatorSeqStoreMock := mock_indexer.NewMockSystemEventCreatorStore(ctrl)

			task := NewSystemEventCreatorTask(testCfg, validatorSeqStoreMock)
			createdSystemEvents, _ := task.getActiveSetPresenceChangeSystemEvents(tt.currHeightList, tt.prevHeightList)

			if len(createdSystemEvents) != tt.expectedCount {
				t.Errorf("unexpected system event count, want %v; got %v", tt.expectedCount, len(createdSystemEvents))
				return
			}

			for i, kind := range tt.expectedKinds {
				if len(createdSystemEvents) > 0 && createdSystemEvents[i].Kind != kind {
					t.Errorf("unexpected system event kind, want %v; got %v", kind, createdSystemEvents[i].Kind)
				}
			}
		})
	}
}

func TestSystemEventCreatorTask_getMissedBlocksSystemEvents(t *testing.T) {
	tests := []struct {
		description           string
		maxValidatorSequences int64
		missedInRowThreshold  int64
		missedForMaxThreshold int64
		prevHeightList        []model.ValidatorSeq
		currHeightList        []model.ValidatorSeq
		lastForValidatorList  [][]model.ValidatorSeq
		errs                  []error
		expectedCount         int
		expectedKinds         []model.SystemEventKind
		expectedErr           error
	}{
		{
			description:           "returns no system events when validator does not have any previous sequences in db",
			maxValidatorSequences: 5,
			missedInRowThreshold:  2,
			missedForMaxThreshold: 2,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{},
			},
			expectedCount: 0,
		},
		{
			description:           "returns no system events when validator does not have any missed blocks in a row",
			maxValidatorSequences: 5,
			missedInRowThreshold:  2,
			missedForMaxThreshold: 2,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns no system events when validator missed 2 blocks in a row",
			maxValidatorSequences: 5,
			missedInRowThreshold:  3,
			missedForMaxThreshold: 5,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns one missed_m_in_row system events when validator missed >= 3 blocks in a row",
			maxValidatorSequences: 5,
			missedInRowThreshold:  3,
			missedForMaxThreshold: 5,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNInRow},
		},
		{
			description:           "returns no missed_m_in_row system events when validator missed >= 3 blocks in a row in the past but current is validated",
			maxValidatorSequences: 5,
			missedInRowThreshold:  3,
			missedForMaxThreshold: 5,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns one missed_m_of_n system events when validator missed 3 blocks",
			maxValidatorSequences: 5,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
		{
			description:           "returns one missed_m_of_n system events when validator missed 3 blocks and max < last list",
			maxValidatorSequences: 3,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
		{
			description:           "returns no missed_m_of_n system events when count of recent not validated > maxValidatorSequences",
			maxValidatorSequences: 5,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns no missed_m_of_n system events when current is validated",
			maxValidatorSequences: 5,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns error when first call to FindLastByAddress fails",
			maxValidatorSequences: 3,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				nil,
			},
			errs:          []error{ErrCouldNotFindByAddress},
			expectedCount: 0,
			expectedErr:   ErrCouldNotFindByAddress,
		},
		{
			description:           "returns error when second call to FindLastByAddress fails",
			maxValidatorSequences: 5,
			missedInRowThreshold:  3,
			missedForMaxThreshold: 5,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
				newValidatorSeq("address1", 1000, 0, false),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
				newValidatorSeq("address1", 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
				},
				nil,
			},
			errs:          []error{nil, ErrCouldNotFindByAddress},
			expectedCount: 0,
			expectedErr:   ErrCouldNotFindByAddress,
		},
		{
			description:           "returns partial system events when second call to FindLastByAddress fails with ErrNotFound",
			maxValidatorSequences: 3,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, true),
				newValidatorSeq("address1", 1000, 0, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, 0, false),
				newValidatorSeq("address1", 1000, 0, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, false),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
					newValidatorSeq(testValidatorAddress, 1000, 0, true),
				},
				nil,
			},
			errs:          []error{nil, store.ErrNotFound},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validatorSeqStoreMock := mock_indexer.NewMockSystemEventCreatorStore(ctrl)

			MaxValidatorSequences = tt.maxValidatorSequences
			MissedInRowThreshold = tt.missedInRowThreshold
			MissedForMaxThreshold = tt.missedForMaxThreshold

			var mockCalls []*gomock.Call
			for i, validatorSeqs := range tt.lastForValidatorList {
				validatorSeq := tt.currHeightList[i]
				if validatorSeq.PrecommitValidated == nil || !*validatorSeq.PrecommitValidated {
					call := validatorSeqStoreMock.EXPECT().FindLastByAddress(gomock.Any(), gomock.Any())

					if len(tt.errs) >= i+1 && tt.errs[i] != nil {
						call = call.Return(nil, tt.errs[i])
					} else {
						call = call.Return(validatorSeqs, nil)
					}

					mockCalls = append(mockCalls, call)
				}
			}
			gomock.InOrder(mockCalls...)

			task := NewSystemEventCreatorTask(testCfg, validatorSeqStoreMock)
			createdSystemEvents, err := task.getMissedBlocksSystemEvents(tt.currHeightList)
			if err == nil && tt.expectedErr != nil {
				t.Errorf("should return error")
				return
			}
			if err != nil && tt.expectedErr != err {
				t.Errorf("unexpected error, want %v; got %v", tt.expectedErr, err)
				return
			}

			if len(createdSystemEvents) != tt.expectedCount {
				t.Errorf("unexpected system event count, want %v; got %v", tt.expectedCount, len(createdSystemEvents))
				return
			}

			for i, kind := range tt.expectedKinds {
				if len(createdSystemEvents) > 0 && createdSystemEvents[i].Kind != kind {
					t.Errorf("unexpected system event kind, want %v; got %v", kind, createdSystemEvents[i].Kind)
				}
			}
		})
	}
}

func testPayload() *payload {
	return &payload{
		Syncable: &model.Syncable{
			Height: testHeight,
			Time:   *types.NewTimeFromTime(time.Now()),
		},
		CurrentHeight: testHeight,
	}
}

func newValidatorSeq(address string, balance int64, commission int64, validated bool) model.ValidatorSeq {
	return model.ValidatorSeq{
		Sequence: &model.Sequence{
			Height: testHeight,
			Time:   *types.NewTimeFromTime(time.Now()),
		},
		Address:             address,
		ActiveEscrowBalance: types.NewQuantityFromInt64(balance),
		Commission:          types.NewQuantityFromInt64(commission),
		PrecommitValidated:  &validated,
	}
}
