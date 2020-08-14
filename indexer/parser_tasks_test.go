package indexer

import (
	"context"
	"reflect"
	"testing"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/types"
)

func TestBlockParserTask_Run(t *testing.T) {
	proposerAddr := "proposerAddr"
	proposerKey := "proposerKey"

	tests := []struct {
		description              string
		rawBlock                 *blockpb.Block
		rawTransactions          []*transactionpb.Transaction
		rawValidators            []*validatorpb.Validator
		expectedTransactionCount int64
		expectedUID              string
	}{
		{"updates empty state",
			nil,
			[]*transactionpb.Transaction{},
			[]*validatorpb.Validator{},
			0,
			"",
		},
		{"updates payload.ParsedBlockData.TransactionsCount",
			nil,
			[]*transactionpb.Transaction{testpbTransaction("t1"), testpbTransaction("t2")},
			[]*validatorpb.Validator{},
			2,
			"",
		},
		{"updates payload.ParsedBlockData.ProposerEntityUID when proposer is in validator list",
			testpbBlock(setBlockProposerAddress(proposerAddr)),
			[]*transactionpb.Transaction{},
			[]*validatorpb.Validator{
				testpbValidator(),
				testpbValidator(setValidatorEntityID(proposerKey), setValidatorAddress(proposerAddr)),
				testpbValidator(),
			},
			0,
			proposerKey,
		},
		{"does not update payload.ParsedBlockData.ProposerEntityUID when proposer is not in validator list",
			testpbBlock(setBlockProposerAddress(proposerAddr)),
			[]*transactionpb.Transaction{},
			[]*validatorpb.Validator{testpbValidator(), testpbValidator()},
			0,
			"",
		},
		{"does not update payload.ParsedBlockData.ProposerEntityUID when there's no block",
			nil,
			[]*transactionpb.Transaction{},
			[]*validatorpb.Validator{testpbValidator(setValidatorEntityID(proposerKey), setValidatorAddress(proposerAddr))},
			0,
			"",
		},
	}

	for _, tt := range tests {
		tt := tt // need to set this since running tests in parallel
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			task := NewBlockParserTask()

			pl := &payload{
				RawBlock:        tt.rawBlock,
				RawValidators:   tt.rawValidators,
				RawTransactions: tt.rawTransactions,
			}

			if err := task.Run(ctx, pl); err != nil {
				t.Errorf("unexpected error on Run, want %v; got %v", nil, err)
				return
			}

			if pl.ParsedBlock.TransactionsCount != tt.expectedTransactionCount {
				t.Errorf("Unexpected ProposerEntityUID, want: %+v, got: %+v", tt.expectedTransactionCount, pl.ParsedBlock.TransactionsCount)
				return
			}

			if pl.ParsedBlock.ProposerEntityUID != tt.expectedUID {
				t.Errorf("Unexpected ProposerEntityUID count, want: %+v, got: %+v", tt.expectedUID, pl.ParsedBlock.ProposerEntityUID)
				return
			}
		})
	}
}

func TestValidatorParserTask_Run(t *testing.T) {
	proposerAddr := "proposerAddr"
	isFalse := false
	isTrue := true
	hundredInBytes := uintToBytes(100, t)

	tests := []struct {
		description                  string
		rawBlock                     *blockpb.Block
		rawStakingState              *statepb.Staking
		rawValidators                []*validatorpb.Validator
		expectedParsedValidatorsData ParsedValidatorsData
	}{
		{"update validator with no block votes",
			testpbBlock(
				setBlockLastCommitVotes(),
				setBlockProposerAddress(proposerAddr),
			),
			nil,
			[]*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
					setTendermintAddress(proposerAddr),
				),
			},
			ParsedValidatorsData{
				"t0": parsedValidator{
					Proposed:             true,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       0,
					TotalShares:          types.NewQuantityFromInt64(0),
				},
			},
		},
		{"updates total shares",
			testpbBlock(
				setBlockLastCommitVotes(),
				setBlockProposerAddress(proposerAddr),
			),
			testpbStaking(
				setStakingDelegationEntry("t0", "entry1", hundredInBytes),
				setStakingDelegationEntry("t0", "entry2", hundredInBytes),
				setStakingDelegationEntry("t1", "entry1", hundredInBytes),
			),
			[]*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
					setTendermintAddress(proposerAddr),
				),
				testpbValidator(
					setValidatorAddress("t1"),
				),
			},
			ParsedValidatorsData{
				"t0": parsedValidator{
					Proposed:             true,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       0,
					TotalShares:          types.NewQuantityFromInt64(200),
				},
				"t1": parsedValidator{
					Proposed:             false,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       1,
					TotalShares:          types.NewQuantityFromInt64(100),
				},
			},
		},
		{"updates PrecommitBlockIdFlag and PrecommitValidated",
			testpbBlock(
				setBlockLastCommitVotes(
					testpbVote(0, 2), // validatorindex=0, blockIDFlag=2
					testpbVote(1, 2),
					testpbVote(2, 1),
				),
				setBlockProposerAddress(proposerAddr),
			),
			nil,
			[]*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
				),
				testpbValidator(
					setValidatorAddress("t1"),
					setTendermintAddress(proposerAddr),
				),
				testpbValidator(
					setValidatorAddress("t2"),
				),
			},
			ParsedValidatorsData{
				"t0": parsedValidator{
					Proposed:             false,
					PrecommitValidated:   &isTrue,
					PrecommitBlockIdFlag: 2,
					PrecommitIndex:       0,
					TotalShares:          types.NewQuantityFromInt64(0),
				},
				"t1": parsedValidator{
					Proposed:             true,
					PrecommitValidated:   &isTrue,
					PrecommitBlockIdFlag: 2,
					PrecommitIndex:       1,
					TotalShares:          types.NewQuantityFromInt64(0),
				},
				"t2": parsedValidator{
					Proposed:             false,
					PrecommitValidated:   &isFalse,
					PrecommitBlockIdFlag: 1,
					PrecommitIndex:       2,
					TotalShares:          types.NewQuantityFromInt64(0),
				},
			},
		},
		{"update validators when there's less votes than validators",
			testpbBlock(
				setBlockLastCommitVotes(
					testpbVote(0, 2),
				),
				setBlockProposerAddress(proposerAddr),
			),
			nil,
			[]*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
					setTendermintAddress(proposerAddr),
				),
				testpbValidator(
					setValidatorAddress("t1"),
				),
			},
			ParsedValidatorsData{
				"t0": parsedValidator{
					Proposed:             true,
					PrecommitValidated:   &isTrue,
					PrecommitBlockIdFlag: 2,
					PrecommitIndex:       0,
					TotalShares:          types.NewQuantityFromInt64(0),
				},
				"t1": parsedValidator{
					Proposed:             false,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       1,
					TotalShares:          types.NewQuantityFromInt64(0),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx := context.Background()

			task := NewValidatorsParserTask()

			pl := &payload{
				RawBlock:        tt.rawBlock,
				RawValidators:   tt.rawValidators,
				RawStakingState: tt.rawStakingState,
			}

			if err := task.Run(ctx, pl); err != nil {
				t.Errorf("unexpected error on Run, want %v; got %v", nil, err)
				return
			}

			if len(pl.ParsedValidators) != len(tt.expectedParsedValidatorsData) {
				t.Errorf("Unexpected number of ParsedValidators, want: %v, got: %v", len(tt.expectedParsedValidatorsData), len(pl.ParsedValidators))
				return
			}

			for key, expected := range tt.expectedParsedValidatorsData {
				val, ok := pl.ParsedValidators[key]
				if !ok {
					t.Errorf("Missing key in payload.ParsedValidators, want: %v", key)
					return
				}
				if !reflect.DeepEqual(val, expected) {
					t.Errorf("Unexpected value in map, want: %+v, got: %+v", expected, val)
				}
			}
		})
	}
}
