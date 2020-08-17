package indexer

import (
	"context"
	"reflect"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/event/eventpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	mock "github.com/figment-networks/oasishub-indexer/mock/client"
	"github.com/golang/mock/gomock"

	"testing"
)

func TestBlockFetcher_Run(t *testing.T) {
	tests := []struct {
		description   string
		expectedBlock *blockpb.Block
		result        error
	}{
		{"returns error if client errors", nil, errTestClient},
		{"updates payload.RawBlock", testpbBlock(), nil},
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
	tests := []struct {
		description   string
		expectedState *statepb.State
		result        error
	}{
		{"returns error if client errors", nil, errTestClient},
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

func TestStakingStateFetcher_Run(t *testing.T) {
	tests := []struct {
		description     string
		expectedStaking *statepb.Staking
		result          error
	}{
		{"returns error if client errors", nil, errTestClient},
		{"updates payload.RawStakingState", testpbStaking(), nil},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockClient := mock.NewMockStateClient(ctrl)
			task := NewStakingStateFetcherTask(mockClient)

			pl := &payload{CurrentHeight: 30}

			mockClient.EXPECT().GetStakingByHeight(pl.CurrentHeight).Return(&statepb.GetStakingByHeightResponse{Staking: tt.expectedStaking}, tt.result).Times(1)

			if result := task.Run(ctx, pl); result != tt.result {
				t.Errorf("want %v; got %v", tt.result, result)
				return
			}

			// skip payload check if there's an error
			if tt.result != nil {
				return
			}

			if !reflect.DeepEqual(pl.RawStakingState, tt.expectedStaking) {
				t.Errorf("want: %+v, got: %+v", tt.expectedStaking, pl.RawStakingState)
				return
			}
		})
	}
}

func TestValidatorFetcher_Run(t *testing.T) {
	tests := []struct {
		description        string
		expectedValidators []*validatorpb.Validator
		result             error
	}{
		{"returns error if client errors", nil, errTestClient},
		{"updates payload.RawValidators", []*validatorpb.Validator{testpbValidator(), testpbValidator()}, nil},
		{"updates payload.RawValidators when client returns empty list", []*validatorpb.Validator{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockClient := mock.NewMockValidatorClient(ctrl)
			task := NewValidatorFetcherTask(mockClient)

			pl := &payload{CurrentHeight: 30}

			mockClient.EXPECT().GetByHeight(pl.CurrentHeight).Return(&validatorpb.GetByHeightResponse{Validators: tt.expectedValidators}, tt.result).Times(1)

			if result := task.Run(ctx, pl); result != tt.result {
				t.Errorf("want %v; got %v", tt.result, result)
				return
			}

			// skip payload check if there's an error
			if tt.result != nil {
				return
			}

			if len(pl.RawValidators) != len(tt.expectedValidators) {
				t.Errorf("wrong number of vallidators in payload; want: %+v, got: %+v", tt.expectedValidators, pl.RawValidators)
				return
			}

			for _, expected := range tt.expectedValidators {
				id := expected.GetNode().GetEntityId()
				for _, validator := range pl.RawValidators {
					if validator.GetNode().GetEntityId() == id && !reflect.DeepEqual(expected, validator) {
						t.Errorf("validators don't match; want: %+v, got: %+v", expected, validator)
					}
				}
			}
		})
	}
}

func TestTransactionFetcher_Run(t *testing.T) {
	tests := []struct {
		description          string
		expectedTransactions []*transactionpb.Transaction
		result               error
	}{
		{"returns error if client errors", nil, errTestClient},
		{"updates payload.RawTransactions", []*transactionpb.Transaction{testpbTransaction("test1"), testpbTransaction("test2")}, nil},
		{"updates payload.RawTransactions when client returns empty list", []*transactionpb.Transaction{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockClient := mock.NewMockTransactionClient(ctrl)
			task := NewTransactionFetcherTask(mockClient)

			pl := &payload{CurrentHeight: 30}

			mockClient.EXPECT().GetByHeight(pl.CurrentHeight).Return(&transactionpb.GetByHeightResponse{Transactions: tt.expectedTransactions}, tt.result).Times(1)

			if result := task.Run(ctx, pl); result != tt.result {
				t.Errorf("want %v; got %v", tt.result, result)
				return
			}

			// skip payload check if there's an error
			if tt.result != nil {
				return
			}

			if len(pl.RawTransactions) != len(tt.expectedTransactions) {
				t.Errorf("wrong number of vallidators in payload; want: %+v, got: %+v", tt.expectedTransactions, pl.RawTransactions)
				return
			}

			for _, expected := range tt.expectedTransactions {
				id := expected.GetPublicKey()
				for _, transaction := range pl.RawTransactions {
					if transaction.GetPublicKey() == id && !reflect.DeepEqual(expected, transaction) {
						t.Errorf("transactions don't match; want: %+v, got: %+v", expected, transaction)
					}
				}
			}
		})
	}
}

func TestEventFetcher_Run(t *testing.T) {
	tests := []struct {
		description    string
		expectedEvents []*eventpb.AddEscrowEvent
		result         error
	}{
		{"returns error if client errors", nil, errTestClient},
		{"updates payload", []*eventpb.AddEscrowEvent{&eventpb.AddEscrowEvent{
			Owner:  "ownerAddr",
			Escrow: "escrowAddr",
			Amount: randBytes(5),
		}}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockClient := mock.NewMockEventClient(ctrl)
			task := NewEventsFetcherTask(mockClient)

			pl := &payload{CurrentHeight: 20}

			mockClient.EXPECT().GetAddEscrowEventsByHeight(pl.CurrentHeight).Return(
				&eventpb.GetAddEscrowEventsByHeightResponse{Events: tt.expectedEvents}, tt.result,
			).Times(1)

			if result := task.Run(ctx, pl); result != tt.result {
				t.Errorf("want %v; got %v", tt.result, result)
				return
			}

			// skip payload check if there's an error
			if tt.result != nil {
				return
			}

			if !reflect.DeepEqual(pl.RawEscrowEvents, tt.expectedEvents) {
				t.Errorf("want: %+v, got: %+v", tt.expectedEvents, pl.RawEscrowEvents)
				return
			}
		})
	}
}
