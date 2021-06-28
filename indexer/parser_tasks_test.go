package indexer

import (
	"context"
	"reflect"
	"testing"

	"github.com/figment-networks/oasis-rpc-proxy/grpc/account/accountpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/block/blockpb"
	"github.com/figment-networks/oasis-rpc-proxy/grpc/debondingdelegation/debondingdelegationpb"
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
		rawTransferEvents            []*eventpb.TransferEvent
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
		{description: "updates validator rewards without commission",
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
					Owner:  "t0",
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
					Owner:  "t1",
					Escrow: "t1",
					Amount: twentyInBytes,
				},
			},
			rawTransferEvents: []*eventpb.TransferEvent{
				{
					From:   commonPoolAddr,
					To:     "t0",
					Amount: twentyInBytes,
				},
				{
					From:   commonPoolAddr,
					To:     "t1",
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
		tt := tt
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
	const currHeight int64 = 20

	t.Run("creates reward and commission balance events", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		task := NewBalanceParserTask()

		pld := &payload{
			CurrentHeight:     currHeight,
			CommonPoolAddress: commonPoolAddr,
			RawValidators: []*validatorpb.Validator{
				testpbValidator(setValidatorAddress(escrowAddr)),
			},
			RawEscrowEvents: &eventpb.EscrowEvents{
				Add: []*eventpb.AddEscrowEvent{
					{
						Owner:     escrowAddr,
						Escrow:    escrowAddr,
						Amount:    uintToBytes(60, t), // event for the automatically escrowed commission reward
						NewShares: uintToBytes(33, t),
					},
					{
						Owner:  commonPoolAddr,
						Escrow: escrowAddr,
						Amount: uintToBytes(240, t), // event for the non-commissioned part of the reward (which only increases existing shares prices)
					},
				},
			},
			RawTransferEvents: []*eventpb.TransferEvent{
				{
					From:   commonPoolAddr,
					To:     escrowAddr,
					Amount: uintToBytes(60, t), // event for the commissioned part of the reward
				},
			},
			RawStakingState: &statepb.Staking{
				Ledger: map[string]*accountpb.Account{
					escrowAddr: {
						Escrow: &accountpb.EscrowAccount{
							Active: &accountpb.SharePool{
								Balance:     uintToBytes(600, t),
								TotalShares: uintToBytes(333, t),
							},
						},
					},
				},
				Delegations: map[string]*delegationpb.DelegationEntry{
					escrowAddr: {
						Entries: map[string]*delegationpb.Delegation{
							escrowAddr:     {Shares: uintToBytes(33, t)},
							delegatorAddr1: {Shares: uintToBytes(100, t)},
							delegatorAddr2: {Shares: uintToBytes(200, t)},
						},
					},
				},
			},
		}

		if err := task.Run(ctx, pld); err != nil {
			t.Errorf("unexpected error on Run, want %v; got %v", nil, err)
			return
		}

		expectedEvents := map[string]model.BalanceEvent{
			delegatorAddr1: {
				Height:        currHeight,
				Address:       delegatorAddr1,
				EscrowAddress: escrowAddr,
				Kind:          model.Reward,
				Amount:        types.NewQuantityFromInt64(80), // account has 100/300 shares, so gets 1/3 * 240 of rewards
			},
			delegatorAddr2: {
				Height:        currHeight,
				Address:       delegatorAddr2,
				EscrowAddress: escrowAddr,
				Kind:          model.Reward,
				Amount:        types.NewQuantityFromInt64(160), // account has 200/300 shares, so gets 2/3 * 240 of rewards
			},
			escrowAddr: {
				Height:        currHeight,
				Address:       escrowAddr,
				EscrowAddress: escrowAddr,
				Kind:          model.Commission,
				Amount:        types.NewQuantityFromInt64(60), // from commission event
			},
		}

		if len(pld.BalanceEvents) != len(expectedEvents) {
			t.Errorf("Unexpected BalanceEvents len, want: %+v, got: %+v", len(expectedEvents), len(pld.BalanceEvents))
			return
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

	t.Run("doesn't create balance event if there's no raw escrow events", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		task := NewBalanceParserTask()

		pld := &payload{
			CurrentHeight:     currHeight,
			CommonPoolAddress: commonPoolAddr,
			RawValidators: []*validatorpb.Validator{
				testpbValidator(setValidatorAddress(escrowAddr)),
			},
			RawEscrowEvents: &eventpb.EscrowEvents{
				Add:  []*eventpb.AddEscrowEvent{},
				Take: []*eventpb.TakeEscrowEvent{},
			},
			RawStakingState: &statepb.Staking{
				Ledger: map[string]*accountpb.Account{
					escrowAddr: {
						Escrow: &accountpb.EscrowAccount{
							Active: &accountpb.SharePool{
								Balance:     uintToBytes(600, t),
								TotalShares: uintToBytes(333, t),
							},
						},
					},
				},
				Delegations: map[string]*delegationpb.DelegationEntry{
					escrowAddr: {
						Entries: map[string]*delegationpb.Delegation{
							escrowAddr:     {Shares: uintToBytes(33, t)},
							delegatorAddr1: {Shares: uintToBytes(100, t)},
							delegatorAddr2: {Shares: uintToBytes(200, t)},
						},
					},
				},
			},
		}

		if err := task.Run(ctx, pld); err != nil {
			t.Errorf("unexpected error on Run, want %v; got %v", nil, err)
			return
		}

		if len(pld.BalanceEvents) != 0 {
			t.Errorf("Unexpected BalanceEvents len, want: %+v, got: %+v", 0, len(pld.BalanceEvents))
			return
		}
	})

	t.Run("creates slash active events", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		task := NewBalanceParserTask()

		pld := &payload{
			CurrentHeight:     currHeight,
			CommonPoolAddress: commonPoolAddr,
			RawValidators: []*validatorpb.Validator{
				testpbValidator(setValidatorAddress(escrowAddr)),
			},
			RawEscrowEvents: &eventpb.EscrowEvents{
				Take: []*eventpb.TakeEscrowEvent{
					{
						Owner:  escrowAddr,
						Amount: uintToBytes(50, t),
					},
				},
			},
			RawStakingState: &statepb.Staking{
				Ledger: map[string]*accountpb.Account{
					escrowAddr: {
						Escrow: &accountpb.EscrowAccount{
							Active: &accountpb.SharePool{
								Balance:     uintToBytes(350, t),
								TotalShares: uintToBytes(1000, t),
							},
						},
					},
				},
				Delegations: map[string]*delegationpb.DelegationEntry{
					escrowAddr: {
						Entries: map[string]*delegationpb.Delegation{
							escrowAddr:     {Shares: uintToBytes(600, t)},
							delegatorAddr1: {Shares: uintToBytes(300, t)},
							delegatorAddr2: {Shares: uintToBytes(100, t)},
						},
					},
				},
			},
		}

		if err := task.Run(ctx, pld); err != nil {
			t.Errorf("unexpected error on Run, want %v; got %v", nil, err)
			return
		}

		expectedEvents := map[string]model.BalanceEvent{
			escrowAddr: {
				Height:        currHeight,
				Address:       escrowAddr,
				EscrowAddress: escrowAddr,
				Kind:          model.SlashActive,
				Amount:        types.NewQuantityFromInt64(30), // 60% of 50
			},
			delegatorAddr1: {
				Height:        currHeight,
				Address:       delegatorAddr1,
				EscrowAddress: escrowAddr,
				Kind:          model.SlashActive,
				Amount:        types.NewQuantityFromInt64(15), // 30% of 50
			},
			delegatorAddr2: {
				Height:        currHeight,
				Address:       delegatorAddr2,
				EscrowAddress: escrowAddr,
				Kind:          model.SlashActive,
				Amount:        types.NewQuantityFromInt64(5), // 10% of 50
			},
		}

		if len(pld.BalanceEvents) != len(expectedEvents) {
			t.Errorf("Unexpected BalanceEvents len, want: %+v, got: %+v", len(expectedEvents), len(pld.BalanceEvents))
			return
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

	t.Run("creates slash debonding events", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		task := NewBalanceParserTask()

		pld := &payload{
			CurrentHeight:     currHeight,
			CommonPoolAddress: commonPoolAddr,
			RawValidators: []*validatorpb.Validator{
				testpbValidator(setValidatorAddress(escrowAddr)),
			},
			RawEscrowEvents: &eventpb.EscrowEvents{
				Take: []*eventpb.TakeEscrowEvent{
					{
						Owner:  escrowAddr,
						Amount: uintToBytes(500, t),
					},
				},
			},
			RawStakingState: &statepb.Staking{
				Ledger: map[string]*accountpb.Account{
					escrowAddr: {
						Escrow: &accountpb.EscrowAccount{
							Debonding: &accountpb.SharePool{
								Balance:     uintToBytes(3500, t),
								TotalShares: uintToBytes(1000, t),
							},
						},
					},
				},
				DebondingDelegations: map[string]*debondingdelegationpb.DebondingDelegationEntry{
					escrowAddr: {
						Entries: map[string]*debondingdelegationpb.DebondingDelegationInnerEntry{
							escrowAddr: {
								DebondingDelegations: []*debondingdelegationpb.DebondingDelegation{
									{Shares: uintToBytes(200, t)},
									{Shares: uintToBytes(150, t)},
									{Shares: uintToBytes(250, t)},
								},
							},
							delegatorAddr1: {
								DebondingDelegations: []*debondingdelegationpb.DebondingDelegation{
									{Shares: uintToBytes(400, t)},
								},
							},
						},
					},
				},
			},
		}

		if err := task.Run(ctx, pld); err != nil {
			t.Errorf("unexpected error on Run, want %v; got %v", nil, err)
			return
		}

		expectedEvents := map[string]model.BalanceEvent{
			escrowAddr: {
				Height:        currHeight,
				Address:       escrowAddr,
				EscrowAddress: escrowAddr,
				Kind:          model.SlashDebonding,
				Amount:        types.NewQuantityFromInt64(300), // 60% of 500
			},
			delegatorAddr1: {
				Height:        currHeight,
				Address:       delegatorAddr1,
				EscrowAddress: escrowAddr,
				Kind:          model.SlashDebonding,
				Amount:        types.NewQuantityFromInt64(200), // 40% of 500
			},
		}

		if len(pld.BalanceEvents) != len(expectedEvents) {
			t.Errorf("Unexpected BalanceEvents len, want: %+v, got: %+v", len(expectedEvents), len(pld.BalanceEvents))
			return
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
