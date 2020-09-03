package indexer

import (
	"context"
	"reflect"
	"testing"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/delegation/delegationpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/event/eventpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/validator/validatorpb"
	"github.com/figment-networks/oasishub-indexer/model"
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
	commonPoolAddr := "commonPoolAddr"
	isFalse := false
	isTrue := true
	hundredInBytes := uintToBytes(100, t)
	twentyInBytes := uintToBytes(20, t)

	tests := []struct {
		description                  string
		rawBlock                     *blockpb.Block
		rawStakingState              *statepb.Staking
		rawValidators                []*validatorpb.Validator
		rawAddEscrowEvents           []*eventpb.AddEscrowEvent
		expectedParsedValidatorsData ParsedValidatorsData
	}{
		{description: "update validator with no block votes",
			rawBlock: testpbBlock(
				setBlockLastCommitVotes(),
				setBlockProposerAddress(proposerAddr),
			),
			rawStakingState: nil,
			rawValidators: []*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
					setTendermintAddress(proposerAddr),
				),
			},
			expectedParsedValidatorsData: ParsedValidatorsData{
				"t0": parsedValidator{
					Proposed:             true,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       0,
					TotalShares:          types.NewQuantityFromInt64(0),
				},
			},
		},
		{description: "updates total shares",
			rawBlock: testpbBlock(
				setBlockLastCommitVotes(),
				setBlockProposerAddress(proposerAddr),
			),
			rawStakingState: testpbStaking(
				setStakingDelegationEntry("t0", "entry1", hundredInBytes),
				setStakingDelegationEntry("t0", "entry2", hundredInBytes),
				setStakingDelegationEntry("t1", "entry1", hundredInBytes),
			),
			rawValidators: []*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
					setTendermintAddress(proposerAddr),
				),
				testpbValidator(
					setValidatorAddress("t1"),
				),
			},
			expectedParsedValidatorsData: ParsedValidatorsData{
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
		{description: "updates PrecommitBlockIdFlag and PrecommitValidated",
			rawBlock: testpbBlock(
				setBlockLastCommitVotes(
					testpbVote(0, 2), // validatorindex=0, blockIDFlag=2
					testpbVote(1, 2),
					testpbVote(2, 1),
				),
				setBlockProposerAddress(proposerAddr),
			),
			rawStakingState: nil,
			rawValidators: []*validatorpb.Validator{
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
			expectedParsedValidatorsData: ParsedValidatorsData{
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
		{description: "update validators when there's less votes than validators",
			rawBlock: testpbBlock(
				setBlockLastCommitVotes(
					testpbVote(0, 2),
				),
				setBlockProposerAddress(proposerAddr),
			),
			rawStakingState: nil,
			rawValidators: []*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
					setTendermintAddress(proposerAddr),
				),
				testpbValidator(
					setValidatorAddress("t1"),
				),
			},
			expectedParsedValidatorsData: ParsedValidatorsData{
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
		{description: "updates validator rewards based on AddEscrowEvents with commonpool owner",
			rawBlock: testpbBlock(
				setBlockLastCommitVotes(),
			),
			rawStakingState: nil,
			rawValidators: []*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
				),
				testpbValidator(
					setValidatorAddress("t1"),
				),
			},
			rawAddEscrowEvents: []*eventpb.AddEscrowEvent{
				{
					Owner:  "not common pool addr",
					Escrow: "t0",
					Amount: hundredInBytes,
				},
				{
					Owner:  commonPoolAddr,
					Escrow: "t1",
					Amount: hundredInBytes,
				},
			},
			expectedParsedValidatorsData: ParsedValidatorsData{
				"t0": parsedValidator{
					Proposed:             false,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       0,
					TotalShares:          types.NewQuantityFromInt64(0),
				},
				"t1": parsedValidator{
					Proposed:             false,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       1,
					TotalShares:          types.NewQuantityFromInt64(0),
					Rewards:              types.NewQuantityFromInt64(100),
				},
			},
		},
		{description: "updates validator rewards based on addescrowevent with the higher amount",
			rawBlock: testpbBlock(
				setBlockLastCommitVotes(),
			),
			rawStakingState: nil,
			rawValidators: []*validatorpb.Validator{
				testpbValidator(
					setValidatorAddress("t0"),
				),
				testpbValidator(
					setValidatorAddress("t1"),
				),
			},
			rawAddEscrowEvents: []*eventpb.AddEscrowEvent{
				{
					Owner:  commonPoolAddr,
					Escrow: "t0",
					Amount: twentyInBytes,
				},
				{
					Owner:  commonPoolAddr,
					Escrow: "t0",
					Amount: hundredInBytes,
				},
				{
					Owner:  commonPoolAddr,
					Escrow: "t1",
					Amount: hundredInBytes,
				},
				{
					Owner:  commonPoolAddr,
					Escrow: "t1",
					Amount: twentyInBytes,
				},
			},
			expectedParsedValidatorsData: ParsedValidatorsData{
				"t0": parsedValidator{
					Proposed:             false,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       0,
					TotalShares:          types.NewQuantityFromInt64(0),
					Rewards:              types.NewQuantityFromInt64(100),
				},
				"t1": parsedValidator{
					Proposed:             false,
					PrecommitValidated:   nil,
					PrecommitBlockIdFlag: 3,
					PrecommitIndex:       1,
					TotalShares:          types.NewQuantityFromInt64(0),
					Rewards:              types.NewQuantityFromInt64(100),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			task := NewValidatorsParserTask()

			pl := &payload{
				RawBlock:          tt.rawBlock,
				RawValidators:     tt.rawValidators,
				RawStakingState:   tt.rawStakingState,
				RawEscrowEvents:   &eventpb.EscrowEvents{Add: tt.rawAddEscrowEvents},
				CommonPoolAddress: commonPoolAddr,
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

func TestBalanceParserTask_Run(t *testing.T) {
	commonPoolAddr := "commonPoolAddr"
	escrowAddr := "escrowAddr"
	delegatorAddr1 := "delegatorAddr1"
	delegatorAddr2 := "delegatorAddr2"

	t.Run("escrowAddr with commission", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		task := NewBalanceParserTask()

		pld := &payload{
			CommonPoolAddress: commonPoolAddr,
			RawValidators: []*validatorpb.Validator{
				testpbValidator(setValidatorAddress(escrowAddr)),
			},
			RawEscrowEvents: &eventpb.EscrowEvents{
				Add: []*eventpb.AddEscrowEvent{
					&eventpb.AddEscrowEvent{
						Owner:  commonPoolAddr,
						Escrow: escrowAddr,
						Amount: uintToBytes(60, t), // commission event
					},
					&eventpb.AddEscrowEvent{
						Owner:  commonPoolAddr,
						Escrow: escrowAddr,
						Amount: uintToBytes(240, t), // reward event
					},
				},
			},
			RawStakingState: &statepb.Staking{
				Ledger: map[string]*accountpb.Account{
					escrowAddr: &accountpb.Account{
						Escrow: &accountpb.EscrowAccount{
							Active: &accountpb.SharePool{
								Balance:     uintToBytes(600, t),
								TotalShares: uintToBytes(333, t),
							},
						},
					},
				},
				Delegations: map[string]*delegationpb.DelegationEntry{
					escrowAddr: &delegationpb.DelegationEntry{
						Entries: map[string]*delegationpb.Delegation{
							// 100% shares are from commission: commission_amount * total_shares / pre_commission_balance = 60 * 300 / 540 = 33
							escrowAddr:     &delegationpb.Delegation{Shares: uintToBytes(33, t)},
							delegatorAddr1: &delegationpb.Delegation{Shares: uintToBytes(100, t)},
							delegatorAddr2: &delegationpb.Delegation{Shares: uintToBytes(200, t)},
						},
					},
				},
			},
		}

		if err := task.Run(ctx, pld); err != nil {
			t.Errorf("unexpected error on Run, want %v; got %v", nil, err)
			return
		}

		if len(pld.BalanceEvents) != 3 {
			t.Errorf("Unexpected BalanceEvents len, want: %+v, got: %+v", 3, len(pld.BalanceEvents))
			return
		}

		expectedEvents := map[string]model.BalanceEvent{
			delegatorAddr1: model.BalanceEvent{
				Height:        0,
				Address:       delegatorAddr1,
				EscrowAddress: escrowAddr,
				Kind:          model.Reward,
				Amount:        types.NewQuantityFromInt64(80), // account has 100/300 shares, so gets 1/3 * 240 of rewards
			},
			delegatorAddr2: model.BalanceEvent{
				Height:        0,
				Address:       delegatorAddr2,
				EscrowAddress: escrowAddr,
				Kind:          model.Reward,
				Amount:        types.NewQuantityFromInt64(160), // account has 200/300 shares, so gets 2/3 * 240 of rewards
			},
			escrowAddr: model.BalanceEvent{
				Height:        0,
				Address:       escrowAddr,
				EscrowAddress: escrowAddr,
				Kind:          model.Commission,
				Amount:        types.NewQuantityFromInt64(60), // from commission event
			},
		}

		for _, event := range pld.BalanceEvents {
			expected, ok := expectedEvents[event.Address]
			if !ok {
				t.Errorf("Unexpected event in payload.BalanceEvents, got: %+v", event)
				return
			}
			if !reflect.DeepEqual(event, expected) {
				t.Errorf("Unexpected event, want: %+v, got: %+v", expected, event)
			}
		}
	})

}
